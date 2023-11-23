package udpclient

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/service/udp"
	"time"
)

func (c *UdpClient) ReadSingleSdo() error {
	return nil
}

func (c *UdpClient) ReadMultiSdo() error {
	return nil
}

func (c *UdpClient) requestId() int {
	c.rid++
	return c.rid
}

func (c *UdpClient) sendBuffSlices(data []byte, packSize int) {
	//	将一包数据切片, 每片1400 byte
	//组成结构: 第一片 : 索引+命令ID+1392(data)
	//后续片 : 索引+命令ID+1400(data)
	//	每生成一片,由生成器推出
	head := data[4:10]
	cmdLen := 10
	cnt := len(data) / packSize
	for i := 0; i < cnt; i++ {
		packIdBytes := udp.GetRawData(int32(i+1), udp.I32)
		message := append(append(append([]byte(StartPackBytes), packIdBytes...), head...), data[packSize*i+cmdLen:packSize*(i+1)+cmdLen]...)
		_, _ = c.udpConnect.Write(message)
		fmt.Printf("%s write % 02x to <%s> from <%s>\n", time.Now(), message, c.udpConnect.RemoteAddr(), c.udpConnect.LocalAddr())
	}
	packIdBytes1 := udp.GetRawData(int32(cnt+1)*(-1), udp.I32)
	message1 := append(append(append([]byte(StartPackBytes), packIdBytes1...), head...), data[packSize*cnt+cmdLen:]...)
	_, _ = c.udpConnect.Write(message1)
	fmt.Printf("%s write % 02x to <%s> from <%s>\n", time.Now(), message1, c.udpConnect.RemoteAddr(), c.udpConnect.LocalAddr())

}

func (c *UdpClient) send(data []byte) error {
	c.sendBuffSlices(data, PackSize)
	return nil
}

func (c *UdpClient) request(data []byte) error {
	_ = c.send(data)
	return nil
}

func (c *UdpClient) connect() {
	err := c.request([]byte("\xff\xff\xff\xff\x01\x00\x00\x00 \x00\x01\x00"))
	if err != nil {
		return
	}
	return
}
func (c *UdpClient) heartBuff() []byte {
	return nil
}
func (c *UdpClient) runHeart() {
	if !c.connected {
		c.connect()
		if c.connected {
			fmt.Println("Udp Is Connected!")
		}
	} else {
		hrtBuff := c.heartBuff()
		c.heardCnt += 1
		ret := c.request(hrtBuff)
		fmt.Println(ret)
		//if !ret.isSuccess {
		//	c.connected = ret.isSuccess
		//}

	}

}
