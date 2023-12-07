package socket_writer

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"time"
)

type Controller interface {
	Read(c net.Conn)
}

type SocketWriter struct {
	Address         string
	KeepAlivePeriod time.Duration

	net.Conn

	Controller
}

func NewSocketWriter(addr string, controller Controller) *SocketWriter {

	return &SocketWriter{
		Address:    addr,
		Controller: controller,
	}
}

func (sw *SocketWriter) Description() string {
	return "Generic socket writer capable of handling multiple socket types."
}

//func (sw *SocketWriter) SampleConfig() string {
//	return `
//  ## URL to connect to
//  # address = "tcp://127.0.0.1:8094"
//  # address = "tcp://example.com:http"
//  # address = "tcp4://127.0.0.1:8094"
//  # address = "tcp6://127.0.0.1:8094"
//  # address = "tcp6://[2001:db8::1]:8094"
//  # address = "udp://127.0.0.1:8094"
//  # address = "udp4://127.0.0.1:8094"
//  # address = "udp6://127.0.0.1:8094"
//  # address = "unix:///tmp/telegraf.sock"
//  # address = "unixgram:///tmp/telegraf.sock"
//
//  ## Optional TLS ControllerConfig
//  # tls_ca = "/etc/telegraf/ca.pem"
//  # tls_cert = "/etc/telegraf/cert.pem"
//  # tls_key = "/etc/telegraf/key.pem"
//  ## Use TLS but skip chain & host verification
//  # insecure_skip_verify = false
//
//  ## Period between keep alive probes.
//  ## Only applies to TCP sockets.
//  ## 0 disables keep alive probes.
//  ## Defaults to the OS configuration.
//  # keep_alive_period = "5m"
//
//  ## Data format to generate.
//  ## Each data format has its own unique set of configuration options, read
//  ## more about them here:
//  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
//  # data_format = "influx"
//`
//}

func (sw *SocketWriter) Connect(timeout time.Duration) error {
	spl := strings.SplitN(sw.Address, "://", 2)
	if len(spl) != 2 {
		return fmt.Errorf("invalid address: %s", sw.Address)
	}

	var c net.Conn
	c, err := net.DialTimeout(spl[0], spl[1], timeout)

	if err != nil {
		return err
	}

	if err := sw.setKeepAlive(c); err != nil {
		log.Printf("unable to configure keep alive (%s): %s", sw.Address, err)
	}

	sw.Conn = c

	go sw.Controller.Read(c)

	return nil
}

func (sw *SocketWriter) setKeepAlive(c net.Conn) error {

	tcpc, ok := c.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("cannot set keep alive on a %s socket", strings.SplitN(sw.Address, "://", 2)[0])
	}
	if sw.KeepAlivePeriod.Nanoseconds() == 0 {
		return tcpc.SetKeepAlive(false)
	}
	if err := tcpc.SetKeepAlive(true); err != nil {
		return err
	}
	return tcpc.SetKeepAlivePeriod(sw.KeepAlivePeriod)
}

// IOWrite writes the given metrics to the destination.
// If an error is encountered, it is up to the caller to retry the same write again later.
// Not parallel safe.
func (sw *SocketWriter) Write(buf []byte) error {

	if sw.Conn == nil {
		return nil
	}

	if _, err := sw.Conn.Write(buf); err != nil {
		//TODO log & keep going with remaining strings
		var _err net.Error
		if ok := errors.Is(err, _err); !ok || !_err.Temporary() {
			// permanent error. close the connection
			sw.Close()
			sw.Conn = nil
			return fmt.Errorf("closing connection: %w", err)
		}
		return err
	}

	return nil
}

// Close closes the connection. Noop if already closed.
func (sw *SocketWriter) Close() error {
	if sw.Conn == nil {
		return nil
	}
	err := sw.Conn.Close()
	sw.Conn = nil
	return err
}
