package udpclient

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/services/udp/udpclient/udps"
	"time"
)

func (c *UdpClient) requestId() uint32 {
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
	rid := udps.Raw2Value(data[4:8], udps.I32).(int32)
	pack := udps.NewFuturePack(rid)
	if _, ok := c.recvBuffMap[rid]; ok {
		c.recvBuffMap[rid].Close()
	}
	c.recvBuffMap[rid] = pack
	//rid := udps.raw2Value(data, SdoDataTypeEnum.U32)
	_ = c.send(data)
	result := c.recvBuffMap[rid].Result(3 * time.Second)
	delete(c.recvBuffMap, rid)
	return result
}

func (c *UdpClient) connect() {
	connBuff := append(append([]byte(udps.SinglePackBytes), udps.GetRawData(c.requestId(), udps.U32)...), []byte(" \x00\x01\x00")...)
	ret := c.request(connBuff)
	if ret.Result == udps.SUCCESS && udps.Raw2Value(ret.Content[len(ret.Content)-2:], udps.U16) == uint16(0x0000) {
		c.connected = true
	}
	return
}
func (c *UdpClient) heartBuff() []byte {
	return append(append(append([]byte(udps.SinglePackBytes), udps.GetRawData(c.requestId(), udps.I32)...), udps.GetRawData(0, udps.U16)...), udps.GetRawData(c.heardCnt, udps.U32)...)
}
func (c *UdpClient) runHeart() {
	if !c.connected {
		c.connect()
		if c.connected {
			fmt.Printf("%s udp connected successfully \n ", time.Now())
		}
	} else {
		hrtBuff := c.heartBuff()
		c.heardCnt += 1
		ret := c.request(hrtBuff)
		if !ret.IsSuccess() {
			fmt.Printf("%s udp connected off \n ", time.Now())
		}
		c.connected = ret.IsSuccess()
	}
}

func (c *UdpClient) ReadSingleSdo(sdoKey string) (error, *udps.Sdo) {
	buff := c.request(udps.SdoPack(udps.GetRawData(sdoKey, udps.String), udps.QSdo, c.requestId()))
	fmt.Println(buff)
	return nil, nil
}

func (c *UdpClient) ReadMultiSdo(sdoKeys []string) (error, map[string]udps.Sdo) {
	return nil, nil
}

func (c *UdpClient) ReadSdoDefines() (error, map[string]udps.Sdo) {
	return nil, nil
}

func (c *UdpClient) ReadSubTopics() (error, map[string]udps.Sdo) {
	return nil, nil
}

func (c *UdpClient) ReadJson() (error, string) {
	return nil, ""
}
func (c *UdpClient) WriteSingleSdo(sdo *udps.Sdo) error {
	return nil
}

func (c *UdpClient) WriteMultiSdo(sdoList []*udps.Sdo) (error, map[string]int) {
	return nil, nil
}
