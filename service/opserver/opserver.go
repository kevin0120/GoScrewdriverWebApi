package opserver

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/service/udp/udpclient"
	"net"
)

type Connect struct {
	start      bool
	tcpConnect net.Conn
	udpClient  *udpclient.UdpClient
}

func StartOpServe(addr string, client *udpclient.UdpClient) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on ", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		op := &Connect{tcpConnect: conn,
			udpClient: client,
		}
		go op.handleConnection(conn)
	}
}

func (c *Connect) handleConnection(conn net.Conn) {
	for {

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			conn.Close()
			return
		}
		request := string(buf[:n])
		fmt.Println("Received request:", request)
		a := c.Deserialize(request)
		if a != nil {
			go a.Handle(c)
		}

	}

}
