package openprotocol

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
	"github.com/kevin0120/GoScrewdriverWebApi/utils"
	"github.com/kevin0120/GoScrewdriverWebApi/utils/socket_writer"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"net"
	"time"
)

const (
	DailTimeout  = 5 * time.Second
	BufferSize   = 65535
	MaxProc      = 4 //共4个协程
	MAX_SEQUENCE = 99999
)

var StrictOpSeq = utils.GetEnvBool("ENV_STRICT_OP_SEQ", false) //是否严格控制OpenProtocol Sequence

func newClientContext(endpoint string, diag Diagnostic, handler IClientHandler, sn string, params *OpenProtocolParams) *clientContext {
	ctx := clientContext{
		needResponse:    atomic.NewBool(false),
		closing:         make(chan struct{}, MaxProc),
		diag:            diag,
		clientHandler:   handler,
		sn:              sn,
		tempResultCurve: tightening_device.NewTighteningCurve(),
		tempPhaseResult: &PhaseResult{},
		params:          params,
	}
	ctx.UpdateStatus(BaseDeviceStatusOffline)
	ctx.sockClient = socket_writer.NewSocketWriter(endpoint, &ctx)
	return &ctx
}

func (c *clientContext) IsNeedResponse() bool {
	return c.needResponse.Load()
}

func (c *clientContext) updateNeedRespFlag(flag bool) {
	c.needResponse.Store(flag)
}

type clientContext struct {
	sn                string
	params            *OpenProtocolParams
	status            atomic.Value
	sockClient        *socket_writer.SocketWriter
	keepAliveCount    atomic.Int32
	keepaliveDeadLine atomic.Value
	sendBuffer        chan []byte
	handlerBuf        chan handlerPkg
	needResponse      *atomic.Bool
	responseChannel   chan *respPkg
	tempResultCurve   *tightening_device.TighteningCurve
	tempPhaseResult   *PhaseResult
	//requestChannel  chan uint32
	sequence   *utils.Sequence
	receiveBuf chan []byte

	closing           chan struct{}
	closinghandleRecv chan struct{}

	diag          Diagnostic
	clientHandler IClientHandler
}

func (c *clientContext) start() {
	c.initProcs()
	go c.connect()
}

func (c *clientContext) stop() {
	for i := 0; i < MaxProc; i++ {
		c.closing <- struct{}{}
	}
}

func (c *clientContext) initProcs() {
	go c.procWrite()
	go c.procHandle()
	go c.procAlive()
	go c.scannerInfoHandle()
}

func (c *clientContext) handlePackageOPPayload(src []byte) error {
	msg := src
	lenMsg := len(msg)

	// 如果头的长度不够
	if lenMsg < LenHeader {
		return errors.New(fmt.Sprintf("OpenProtocl SeqNumber:%s Head Is Error: %s", c.sn, string(msg)))
	}

	var header OpenProtocolHeader
	header.Deserialize(string(msg[0:LenHeader]))

	// 如果body的长度匹配
	if header.LEN == lenMsg-LenHeader {
		pkg := handlerPkg{
			SN:     c.sn,
			Header: header,
			Body:   string(msg[LenHeader : LenHeader+header.LEN]),
		}

		c.handlerBuf <- pkg
	} else {
		return errors.New(fmt.Sprintf("OpenProtocol SeqNumber:%s Body Len Err: %s", c.sn, hex.EncodeToString(msg)))
	}

	return nil
}

func (c *clientContext) procHandleRecv() {
	c.receiveBuf = make(chan []byte, BufferSize)

	for {
		select {
		case buf := <-c.receiveBuf:
			err := c.handlePackageOPPayload(buf)
			if err != nil {
				//数据需要丢弃
				c.diag.Error("handlePackageOPPayload Error", err)
				c.diag.Debug(fmt.Sprintf("procHandleRecv Raw Msg:%s", string(buf)))
			}

		case <-c.closinghandleRecv:
			c.diag.Debug("procHandleRecv Exit")
			return
		}

	}
}

func (c *clientContext) procWrite() {
	c.sendBuffer = make(chan []byte, BufferSize)

	for {
		select {
		case v := <-c.sendBuffer:
			err := c.sockClient.Write(v)
			if err != nil {
				c.diag.Error("IOWrite Data Fail", err)
			} else {
				c.diag.Debug(fmt.Sprintf("OpenProtocol Send %s: %s", c.sn, string(v)))
				c.updateKeepAliveDeadLine()
			}

			time.Sleep(100 * time.Millisecond)

		case <-c.closing:
			c.diag.Debug("procWrite Exit")
			return
		}
	}
}

func (ctx *clientContext) scannerInfoHandle() {

	//apply := func(_ context.Context, i interface{}) (interface{}, error) {
	//	return i, nil
	//}
	//
	//ctx.vinSubscribeBuf = make(chan rxgo.Item, 4)
	//observable := rxgo.FromChannel(ctx.vinSubscribeBuf).DistinctUntilChanged(apply)
	//ch := observable.Observe()
	//for {
	//	select {
	//	case <-ctx.closing:
	//		ctx.diag.Debug("scannerInfoHandle Exit")
	//		return
	//
	//	case vinInfo := <-ch:
	//		c := ctx.clientHandler
	//		bc := vinInfo.V.(string)
	//		ctx.diag.Info(fmt.Sprintf("从控制器收到条码信息: %s", bc))
	//		ss := scanner.ScannerRead{
	//			Src:     tightening_device.TIGHTENING_DEVICE_TYPE_CONTROLLER,
	//			SN:      c.SerialNumber(),
	//			Barcode: bc,
	//		}
	//		c.doDispatch(dispatcherbus.DispatcherScannerData, ss)
	//	}
	//}
}

func (ctx *clientContext) procHandle() {
	ctx.handlerBuf = make(chan handlerPkg, BufferSize)

	for {
		select {
		case pkg := <-ctx.handlerBuf:
			err := ctx.clientHandler.handleMsg(&pkg, ctx)
			if err != nil {
				ctx.diag.Error("Open IProtocol handleMsg Fail", err)
			}

		case <-ctx.closing:
			ctx.diag.Debug("procHandle Exit")
			return
		}
	}
}

func (c *clientContext) procAlive() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if c.Status() == BaseDeviceStatusOffline {
				continue
			}

			if c.KeepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				c.sendKeepalive()
				c.updateKeepAliveDeadLine() //更新keepalivedeadline
				c.addKeepAliveCount()
			}

		case <-c.closing:
			c.diag.Debug("procAlive Exit")
			return

		}
	}
}

//func (c *clientContext) Read(conn net.Conn) {
//	defer func() {
//		if err := conn.Close(); err != nil {
//			c.diag.Error("Client Close Error ", err)
//		}
//
//		c.closinghandleRecv <- struct{}{}
//	}()
//
//	buf := make([]byte, BufferSize)
//	//var err error
//
//	for {
//		//if err = conn.SetReadDeadline(time.Now().Add(c.params.KeepAlivePeriod * time.Duration(c.params.MaxKeepAliveCheck)).Add(1 * time.Second)); err != nil {
//		//	c.diag.Error("SetReadDeadline Failed ", err)
//		//	break
//		//}
//		n, err1 := conn.Read(buf)
//
//		fmt.Println("接收到的原始数据为:", string(buf[:n]), "结束")
//		if err1 != nil {
//			break
//		}
//		dst := make([]byte, n-1)
//		copy(dst, buf[:n-1])
//		c.updateKeepAliveCount(0)
//
//		c.receiveBuf <- dst
//
//	}
//
//	c.handleStatus(BaseDeviceStatusOffline)
//	c.clientHandler.HandleStatus(c.sn, BaseDeviceStatusOffline)
//	c.clientHandler.UpdateToolStatus(c.sn, BaseDeviceStatusOffline)
//}

func (c *clientContext) Read(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			c.diag.Error("Client Close Error ", err)
		}

		c.closinghandleRecv <- struct{}{}
	}()

	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		s := byte(OpTerminal)
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, s); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	buf := make([]byte, BufferSize)

	newScanner := bufio.NewScanner(conn)
	newScanner.Buffer(buf, BufferSize)

	newScanner.Split(splitFunc)

	var err error

	for newScanner.Scan() {
		if err = conn.SetReadDeadline(time.Now().Add(c.params.KeepAlivePeriod * time.Duration(c.params.MaxKeepAliveCheck)).Add(1 * time.Second)); err != nil {
			c.diag.Error("SetReadDeadline Failed ", err)
			break
		}

		b := newScanner.Bytes()

		dst := make([]byte, len(b))
		copy(dst, b)
		c.updateKeepAliveCount(0)
		//fmt.Println(fmt.Sprintf("接收到的原始数据为: %s结束", dst))
		c.receiveBuf <- dst
	}

	c.handleStatus(BaseDeviceStatusOffline)
	c.clientHandler.HandleStatus(c.sn, BaseDeviceStatusOffline)
	c.clientHandler.UpdateToolStatus(c.sn, BaseDeviceStatusOffline)
}

func (c *clientContext) Status() string {
	return c.status.Load().(string)
}

func (c *clientContext) UpdateStatus(status string) {
	c.status.Store(status)
}

func (c *clientContext) handleStatus(status string) {

	if status != c.Status() {

		c.UpdateStatus(status)

		if status == BaseDeviceStatusOffline {

			// 断线重连
			go c.connect()
		}
	}
}

func (c *clientContext) KeepAliveCount() int32 {
	return c.keepAliveCount.Load()
}

func (c *clientContext) updateKeepAliveCount(i int32) {
	c.keepAliveCount.Swap(i)
}

func (c *clientContext) addKeepAliveCount() {
	c.keepAliveCount.Inc()
}

func (c *clientContext) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(c.params.KeepAlivePeriod))
}

func (c *clientContext) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *clientContext) sendKeepalive() {
	if c.Status() == BaseDeviceStatusOffline {
		return
	}

	keepAlive := GeneratePackage(MID_9999_ALIVE, DefaultRev, true, "", "", "")
	c.Write([]byte(keepAlive))
}

func (c *clientContext) Write(buf []byte) {
	c.sendBuffer <- buf
}

func (c *clientContext) resetConn() {
	c.UpdateStatus(BaseDeviceStatusOffline)

	//c.requestChannel = make(chan uint32, 1024)
	c.sequence = utils.CreateSequence(MAX_SEQUENCE)
	c.responseChannel = make(chan *respPkg)
}

func (c *clientContext) connect() {
	c.diag.Debug(fmt.Sprintf("OpenProtocol SeqNumber:%s Connecting ...", c.sn))
	c.updateKeepAliveDeadLine()
	c.resetConn()

	for {
		err := c.sockClient.Connect(DailTimeout)
		if err != nil {
			c.diag.Error("connect", err)
		} else {
			c.diag.Debug(fmt.Sprintf("OpenProtocol SeqNumber:%s Connected", c.sn))
			break
		}

		time.Sleep(1 * time.Second)
	}

	c.handleStatus(BaseDeviceStatusOnline)
	c.clientHandler.HandleStatus(c.sn, BaseDeviceStatusOnline)
	c.clientHandler.UpdateToolStatus(c.sn, BaseDeviceStatusOnline)
	c.closinghandleRecv = make(chan struct{}, 1)
	go c.procHandleRecv()

	time.Sleep(100 * time.Millisecond)

	operation := func() error {

		err := c.startComm()
		if err != nil {
			return err
		}
		return nil
	}
	err := backoff.RetryNotify(operation, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		c.diag.Debug(fmt.Sprintf("start Cmd send Failed: %s, retry after %v", err.Error(), duration))
	})
	if err != nil {
		c.diag.Error(fmt.Sprintf("Start Comm Failed: %s", c.sn), err)
	}
}

func (c *clientContext) startComm() error {
	reply, err := c.ProcessRequest(MID_0001_START, false, "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != requestErrors["00"] && reply.(string) != requestErrors["96"] {
		return errors.New(reply.(string))
	}

	return nil
}

func (c *clientContext) SendOpenProtocolAckMsg(mid string, station string, spindle string, data string) (err error) {
	rev, err := c.clientHandler.GetVendorMid(mid)
	if err != nil {
		return
	}

	if c.Status() == BaseDeviceStatusOffline {
		return errors.New(BaseDeviceStatusOffline)
	}

	msg := GeneratePackage(mid, rev, true, station, spindle, data) // no ack always true
	message := []byte(msg)
	c.Write(message)
	return
}

func (c *clientContext) sendOpenProtocolRequest(message []byte, sequence uint32) (*respPkg, error) {
	if sequence == 0 || message == nil {
		return nil, errors.New("sendOpenProtocolRequest, Validate Error")
	}
	c.Write(message)
	ctx, cancel := context.WithTimeout(context.Background(), c.params.MaxReplyTime)
	defer cancel()
	c.updateNeedRespFlag(true)
	for {
		select {
		case <-ctx.Done():
			c.updateNeedRespFlag(false)
			fmt.Println(string(message))
			return nil, errors.Wrap(ctx.Err(), "sendOpenProtocolRequest11")
		case pkg := <-c.responseChannel:
			if !StrictOpSeq {
				return pkg, nil
			}
			if pkg.Seq > sequence {
				err := errors.New("Sequence Is Not Equal, Actual Greater Than Expect!")
				c.diag.Error("", err)
				c.updateNeedRespFlag(false)
				return pkg, err
			}
			if pkg.Seq != sequence {
				c.diag.Error("", errors.Errorf("Sequence Is Not Equal, Expect: %d, Actual:%d", sequence, pkg.Seq))
			} else {
				c.updateNeedRespFlag(false)
				return pkg, nil
			}
		}
	}

}

func (c *clientContext) ProcessRequest(mid string, noack bool, station string, spindle string, data string) (interface{}, error) {
	rev, err := c.clientHandler.GetVendorMid(mid)
	if err != nil {
		return nil, err
	}

	if c.Status() == BaseDeviceStatusOffline {
		return nil, errors.New(BaseDeviceStatusOffline)
	}

	msg := GeneratePackage(mid, rev, noack, station, spindle, data)
	message := []byte(msg)
	//if noack {
	//	c.Write(message)
	//	return requestErrors["00"], nil
	//}

	seq := c.sequence.GetSequence()
	if pkg, err := c.sendOpenProtocolRequest(message, seq); err != nil {
		return "", err
	} else {
		return pkg.Body, err
	}

}
