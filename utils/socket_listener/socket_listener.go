package socket_listener

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Protocol interface {
	Parse(msg []byte) error
	Read(c net.Conn)
	NewConn(c net.Conn)
}

type setReadBufferer interface {
	SetReadBuffer(bytes int) error
}

type streamSocketListener struct {
	net.Listener
	*SocketListener

	sockType string
	Port     string

	connections    map[string]net.Conn
	connectionsMtx sync.Mutex
}

func (ssl *streamSocketListener) listen() {

	for {
		l, err := net.Listen("tcp", ssl.Port)
		if err == nil {
			ssl.Listener = l
			break
		} else {
			log.Printf("listen err:%s", err)
		}

		time.Sleep(300 * time.Millisecond)
	}

	ssl.connections = map[string]net.Conn{}

	for {
		c, err := ssl.Accept()
		if err != nil {
			if !strings.HasSuffix(err.Error(), ": use of closed network connection") {
				log.Printf("streamSocket accept fail %s", err)
			}
			break
		}

		if ssl.ReadBufferSize > 0 {
			if srb, ok := c.(setReadBufferer); ok {
				if err := srb.SetReadBuffer(ssl.ReadBufferSize); err != nil {
					log.Printf("SetReadBuffer error: %v", err)
				}
			} else {
				log.Printf("W! Unable to set read buffer on a %s socket", ssl.sockType)
			}
		}

		ssl.connectionsMtx.Lock()
		if ssl.MaxConnections > 0 && len(ssl.connections) >= ssl.MaxConnections {
			ssl.connectionsMtx.Unlock()
			c.Close()
			continue
		}
		ssl.connections[c.RemoteAddr().String()] = c
		ssl.connectionsMtx.Unlock()

		ssl.NewConn(c)

		if err := ssl.setKeepAlive(c); err != nil {
			log.Printf("unable to configure keep alive (%s): %s", ssl.ServiceAddress, err)
		}

		go ssl.Read(c)
	}

	ssl.connectionsMtx.Lock()
	for _, c := range ssl.connections {
		c.Close()
	}
	ssl.connectionsMtx.Unlock()
}

func (ssl *streamSocketListener) setKeepAlive(c net.Conn) error {

	tcpc, ok := c.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("cannot set keep alive on a %s socket", strings.SplitN(ssl.ServiceAddress, "://", 2)[0])
	}
	if ssl.KeepAlivePeriod.Nanoseconds() == 0 {
		return tcpc.SetKeepAlive(false)
	}
	if err := tcpc.SetKeepAlive(true); err != nil {
		return err
	}
	return tcpc.SetKeepAlivePeriod(ssl.KeepAlivePeriod)
}

func (ssl *streamSocketListener) RemoveConnection(c net.Conn) {
	ssl.connectionsMtx.Lock()
	delete(ssl.connections, c.RemoteAddr().String())
	ssl.connectionsMtx.Unlock()
}

type packetSocketListener struct {
	net.PacketConn
	*SocketListener
}

func (psl *packetSocketListener) listen() {
	buf := make([]byte, 64*1024) // 64kb - maximum size of IP packet
	for {
		n, _, err := psl.ReadFrom(buf)
		if err != nil {
			if !strings.HasSuffix(err.Error(), ": use of closed network connection") {
				log.Printf("UDP read error %s", err)
			}
			break
		}

		err = psl.Parse(buf[:n])
		if err != nil {
			log.Printf("unable to parse incoming packet: %s", err)
			//TODO rate limit
			continue
		}
	}
}

func (psl *packetSocketListener) RemoveConnection(c net.Conn) {
}

type InterListener interface {
	RemoveConnection(c net.Conn)
	Close() error
}

type SocketListener struct {
	ServiceAddress  string
	MaxConnections  int
	ReadBufferSize  int
	ReadTimeout     time.Duration
	KeepAlivePeriod time.Duration

	Protocol //协议解析器服务

	InterListener
}

func (sl *SocketListener) Description() string {
	return "Generic socket listener capable of handling multiple socket types."
}

//func (sl *SocketListener) SampleConfig() string {
//	return `
//  ## URL to listen on
//  # service_address = "tcp://:8094"
//  # service_address = "tcp://127.0.0.1:http"
//  # service_address = "tcp4://:8094"
//  # service_address = "tcp6://:8094"
//  # service_address = "tcp6://[2001:db8::1]:8094"
//  # service_address = "udp://:8094"
//  # service_address = "udp4://:8094"
//  # service_address = "udp6://:8094"
//  # service_address = "unix:///tmp/telegraf.sock"
//  # service_address = "unixgram:///tmp/telegraf.sock"
//
//  ## Maximum number of concurrent connections.
//  ## Only applies to stream sockets (e.g. TCP).
//  ## 0 (default) is unlimited.
//  # max_connections = 1024
//
//  ## IORead timeout.
//  ## Only applies to stream sockets (e.g. TCP).
//  ## 0 (default) is unlimited.
//  # read_timeout = "30s"
//
//  ## Optional TLS configuration.
//  ## Only applies to stream sockets (e.g. TCP).
//  # tls_cert = "/etc/telegraf/cert.pem"
//  # tls_key  = "/etc/telegraf/key.pem"
//  ## Enables client authentication if set.
//  # tls_allowed_cacerts = ["/etc/telegraf/clientca.pem"]
//
//  ## Maximum socket buffer size in bytes.
//  ## For stream sockets, once the buffer fills up, the sender will start backing up.
//  ## For datagram sockets, once the buffer fills up, metrics will start dropping.
//  ## Defaults to the OS default.
//  # read_buffer_size = 65535
//
//  ## Period between keep alive probes.
//  ## Only applies to TCP sockets.
//  ## 0 disables keep alive probes.
//  ## Defaults to the OS configuration.
//  # keep_alive_period = "5m"
//
//  ## Data format to consume.
//  ## Each data format has its own unique set of configuration options, read
//  ## more about them here:
//  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
//  # data_format = "influx"
//`
//}

func (sl *SocketListener) Gather() error {
	return nil
}

func (sl *SocketListener) Start() error {
	spl := strings.SplitN(sl.ServiceAddress, "://", 2)
	if len(spl) != 2 {
		return fmt.Errorf("invalid services address: %s", sl.ServiceAddress)
	}

	if spl[0] == "unix" || spl[0] == "unixpacket" || spl[0] == "unixgram" {
		// no good way of testing for "file does not exist".
		// Instead just ignore error and blow up when we try to listen, which will
		// indicate "address already in use" if file existed and we couldn't remove.
		os.Remove(spl[1])
	}

	switch spl[0] {
	case "tcp", "tcp4", "tcp6", "unix", "unixpacket":
		var (
			err error
			l   net.Listener
		)

		if err != nil {
			return nil
		}

		ssl := &streamSocketListener{
			Listener:       l,
			SocketListener: sl,
			sockType:       spl[0],
			Port:           spl[1],
		}

		sl.InterListener = ssl
		go ssl.listen()
	case "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unixgram":
		pc, err := net.ListenPacket(spl[0], spl[1])
		if err != nil {
			return err
		}

		if sl.ReadBufferSize > 0 {
			if srb, ok := pc.(setReadBufferer); ok {
				if err := srb.SetReadBuffer(sl.ReadBufferSize); err != nil {
					log.Printf("SetReadBuffer error: %v", err)
				}
			} else {
				log.Printf("W! Unable to set read buffer on a %s socket", spl[0])
			}
		}

		psl := &packetSocketListener{
			PacketConn:     pc,
			SocketListener: sl,
		}

		sl.InterListener = psl
		go psl.listen()
	default:
		return fmt.Errorf("unknown protocol '%s' in '%s'", spl[0], sl.ServiceAddress)
	}

	if spl[0] == "unix" || spl[0] == "unixpacket" || spl[0] == "unixgram" {
		sl.InterListener = unixCloser{path: spl[1], closer: sl.InterListener}
	}

	return nil
}

func (sl *SocketListener) Stop() {
	if sl.InterListener != nil {
		sl.Close()
		sl.InterListener = nil
	}
}

func NewSocketListener(addr string, protocol Protocol, readBufSize int, maxConnections int) *SocketListener {

	return &SocketListener{
		ServiceAddress: addr,
		ReadBufferSize: readBufSize,
		MaxConnections: maxConnections,
		//KeepAlivePeriod: time.Second * 3, //默认keepalive 周期３秒
		Protocol: protocol,
	}
}

type unixCloser struct {
	path   string
	closer InterListener
}

func (uc unixCloser) Close() error {
	err := uc.closer.Close()
	os.Remove(uc.path) // ignore error
	return err
}

func (uc unixCloser) RemoveConnection(c net.Conn) {
}
