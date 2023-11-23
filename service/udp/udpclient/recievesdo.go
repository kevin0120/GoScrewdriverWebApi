package udpclient

import (
	"encoding/hex"
	"fmt"
	"time"
)

func (c *UdpClient) HandleReceive(data []byte) {
	fmt.Printf("%s read %s from <%s> to <%s>\n", time.Now(), hex.EncodeToString(data), c.udpConnect.RemoteAddr(), c.udpConnect.LocalAddr())
}
