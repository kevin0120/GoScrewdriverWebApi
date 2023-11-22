package udpclient

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/config"
	"net"
	"time"
)

type UdpClient struct {
	udpConnect *net.UDPConn
	rid        int
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewClient() *UdpClient {
	conf := config.GetConfig()
	// 解析服务器地址
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(conf.UdpClient.RemoteHost),
		Port: conf.UdpClient.RemotePort,
	}
	// 解析本地地址
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: conf.UdpClient.LocalPort,
	}
	//创建连接
	conn, err := net.DialUDP("udp", laddr, udpAddr)
	if err != nil {
		fmt.Println("DialUDP error:", err)
	}

	client := UdpClient{
		rid:        0,
		udpConnect: conn,
	}
	// 创建一个取消上下文和等待组
	client.ctx, client.cancel = context.WithCancel(context.Background())

	//go client.Heart()
	return &client
}

func (c *UdpClient) Open() {
	data := make([]byte, 1024)
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				fmt.Println("接收到退出信号")
				return
			default:
				fmt.Println("协程运行中...")
				c.udpConnect.SetReadDeadline(time.Now().Add(3 * time.Second))
				n, err := c.udpConnect.Read(data)
				if err == nil {
					fmt.Printf("read %s from <%s>\n", hex.EncodeToString(data[:n]), c.udpConnect.RemoteAddr())
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
				fmt.Println("协程运行中2...")
				c.RequestAsync()
				time.Sleep(3 * time.Second)
			}

		}
	}()

}
func (c *UdpClient) Close() {
	c.cancel()
}

func (c *UdpClient) requestId() int {
	c.rid++
	return c.rid
}

func (c *UdpClient) HandleReceive() {

}

func (c *UdpClient) request(data []byte) error {
	_, err := c.udpConnect.Write(data)
	if err != nil {
		return err
	}
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
	//fmt.Println(a)
	_, err := c.udpConnect.Write([]byte("\xff\xff\xff\xff\xff\xff\x01\x00\x00\x00 \x00\x01\x00"))
	if err != nil {
		return err
	}
	return nil
}

func (c *UdpClient) RequestMulti() error {
	return nil
}
