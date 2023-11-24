package udpclient

import (
	"encoding/binary"
	"fmt"
	"time"
)

func (c *UdpClient) HandleReceive(data []byte) {
	fmt.Printf("%s read % 02x from <%s> to <%s>\n", time.Now(), data, c.udpConnect.RemoteAddr(), c.udpConnect.LocalAddr())
	pid := binary.LittleEndian.Uint32(data[2:6])
	fmt.Println(pid)
	rid := int(binary.LittleEndian.Uint32(data[6:10]))
	if _, ok := c.recvBuffMap[rid]; ok {
		if c.recvBuffMap[rid].Head == "" {

		}
		c.recvBuffMap[rid].TimeOut()
	}
}
