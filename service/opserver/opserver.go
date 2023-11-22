package opserver

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/config"
	"github.com/kevin0120/GoScrewdriverWebApi/service/udp/udpclient"
	"net"
	"os"
	"strconv"
)

type Connect struct {
	start      bool
	tcpConnect net.Conn
	udpClient  *udpclient.UdpClient
}

func StartOpServe() {

	conf := config.GetConfig()
	// 获取命令行输入的参数
	// 检查是否至少有一个参数传入
	port := conf.OpPort
	if len(os.Args) >= 2 {
		port, _ = strconv.Atoi(os.Args[1])
	}
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on ", addr)

	udp := udpclient.NewClient()
	udp.Open()

	for {
		conn, err := listener.Accept()
		udp.Close()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		op := &Connect{tcpConnect: conn,
			udpClient: udp,
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
