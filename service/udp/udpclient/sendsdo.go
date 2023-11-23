package udpclient

import "fmt"

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
	//_, err := c.udpConnect.Write(data)
	//if err != nil {
	//	return err
	//}
	//return nil
}

func (c *UdpClient) send(data []byte) error {
	_, err := c.udpConnect.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *UdpClient) request(data []byte) error {
	_ = c.send(data)
	return nil
}

func (c *UdpClient) connect() {
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
