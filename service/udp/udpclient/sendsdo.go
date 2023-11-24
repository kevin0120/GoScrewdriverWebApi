package udpclient

import (
	"encoding/binary"
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/service/udp/udpclient/udps"
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
		packIdBytes := udps.GetRawData(int32(i+1), udps.I32)
		message := append(append(append([]byte(udps.StartPackBytes), packIdBytes...), head...), data[packSize*i+cmdLen:packSize*(i+1)+cmdLen]...)
		_, _ = c.udpConnect.Write(message)
		fmt.Printf("%s write % 02x to <%s> from <%s>\n", time.Now(), message, c.udpConnect.RemoteAddr(), c.udpConnect.LocalAddr())
	}
	packIdBytes1 := udps.GetRawData(int32(cnt+1)*(-1), udps.I32)
	message1 := append(append(append([]byte(udps.StartPackBytes), packIdBytes1...), head...), data[packSize*cnt+cmdLen:]...)
	_, _ = c.udpConnect.Write(message1)
	fmt.Printf("%s write % 02x to <%s> from <%s>\n", time.Now(), message1, c.udpConnect.RemoteAddr(), c.udpConnect.LocalAddr())

}

func (c *UdpClient) send(data []byte) error {
	c.sendBuffSlices(data, udps.PackSize)
	return nil
}

func (c *UdpClient) request(data []byte) *udps.FutureData {
	rid := int(binary.LittleEndian.Uint32(data[4:8]))
	pack := udps.NewFuturePack(rid)
	if _, ok := c.recvBuffMap[rid]; ok {
		c.recvBuffMap[rid].TimeOut()
	}
	c.recvBuffMap[rid] = pack
	//rid := udps.raw2Value(data, SdoDataTypeEnum.U32)
	_ = c.send(data)
	result := c.recvBuffMap[rid].Result(3 * time.Second)
	delete(c.recvBuffMap, rid)
	return result
}

func (c *UdpClient) connect() {
	_ = c.request([]byte("\xff\xff\xff\xff\x01\x00\x00\x00 \x00\x01\x00"))
	//if ret.result == SUCCESS and ret.content[-2:] == ""\x00\x00':
	//FILE_LOG.debug(f'udp send connection buff success !!!')
	//self.connected = True
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
