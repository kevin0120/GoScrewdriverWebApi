package udpclient

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/service/udp/udpclient/udps"
	"time"
)

func (c *UdpClient) HandleReceive(data []byte) {
	fmt.Printf("%s read % 02x from <%s> to <%s>\n", time.Now(), data, c.udpConnect.RemoteAddr(), c.udpConnect.LocalAddr())
	pid := udps.Raw2Value(data[2:6], udps.I32).(int32)
	rid := udps.Raw2Value(data[6:10], udps.I32).(int32)
	if _, ok := c.recvBuffMap[rid]; ok {
		if len(c.recvBuffMap[rid].Head) == 0 {
			c.recvBuffMap[rid].Head = data[2:12]
		}
		c.recvBuffMap[rid].Add(pid, data[12:])
		//c.recvBuffMap[rid].TimeOut()
	}
}
