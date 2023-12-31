package udpclient

import (
	"context"
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/services/udp/udpclient/udps"
	"net"
	"time"
)

type UdpClient struct {
	udpConnect  *net.UDPConn
	connected   bool
	heardCnt    uint32
	recvBuffMap map[int32]*udps.SyncUdpFuturePack
	rid         uint32
	timeoutMs   int
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewClient(timeoutMs int) *UdpClient {
	client := UdpClient{
		timeoutMs:   timeoutMs,
		recvBuffMap: map[int32]*udps.SyncUdpFuturePack{},
		rid:         0,
		connected:   false,
		heardCnt:    1,
	}
	//go client.Heart()
	return &client
}

func (c *UdpClient) Open() {
	data := make([]byte, 2048)
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				fmt.Println("接收到退出信号")
				return
			default:
				c.udpConnect.SetReadDeadline(time.Now().Add(3 * time.Second))
				n, err := c.udpConnect.Read(data)
				if err == nil {
					c.HandleReceive(data[:n])
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				fmt.Println("接收到退出信号2")
				return
			default:
				go c.runHeart()
				time.Sleep(3 * time.Second)
			}
		}
	}()
}
func (c *UdpClient) Close() {
	c.cancel()
}

func (c *UdpClient) ConnectToServer(sIpAddr string, sdoPort int, topicPort int) {
	// 解析服务器地址
	c.HandleHttpRequest()
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(sIpAddr),
		Port: sdoPort,
	}
	// 解析本地地址
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 50005,
	}
	//创建连接
	conn, err := net.DialUDP("udp", laddr, udpAddr)
	if err != nil {
		fmt.Println("DialUDP error:", err)
	}

	c.udpConnect = conn
	// 创建一个取消上下文和等待组
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.Open()
	return
}

func (c *UdpClient) ReadMultiSdoCircle(quarySdo []string) error {
	return nil
}

func (c *UdpClient) RequestAsync() error {
	//a := []byte("\xff\xff\xff\xff\x01\x00\x00\x00\x00\x01\x000x300005")
	//fmt.Println(hex.EncodeToString(a))
	//var a []byte
	//a = SdoDataTypeEnum.GetRawData("path", SdoDataTypeEnum.String)
	//a = SdoDataTypeEnum.GetRawData(true,
	//	SdoDataTypeEnum.Bool)
	//a = SdoDataTypeEnum.GetRawData(int8(-3),
	//	SdoDataTypeEnum.I8)
	//a = SdoDataTypeEnum.GetRawData(uint8(3),
	//	SdoDataTypeEnum.U8)
	//a = SdoDataTypeEnum.GetRawData(int16(-3),
	//	SdoDataTypeEnum.I16)
	//a = SdoDataTypeEnum.GetRawData(uint16(3),
	//	SdoDataTypeEnum.U16)
	//a = SdoDataTypeEnum.GetRawData(int32(-3),
	//	SdoDataTypeEnum.I32)
	//a = SdoDataTypeEnum.GetRawData(uint32(3),
	//	SdoDataTypeEnum.U32)
	//
	//a = SdoDataTypeEnum.GetRawData(int64(-3),
	//	SdoDataTypeEnum.I64)
	//a = SdoDataTypeEnum.GetRawData(uint64(3),
	//	SdoDataTypeEnum.U64)
	//a = SdoDataTypeEnum.GetRawData(float32(0.3),
	//	SdoDataTypeEnum.F32)
	//a = SdoDataTypeEnum.GetRawData(0.3,
	//	SdoDataTypeEnum.F64)

	return nil
}

func (c *UdpClient) RequestMulti() error {
	return nil
}
