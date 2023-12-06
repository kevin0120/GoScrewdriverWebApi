package openprotocol

import (
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/tightening_device"
)

func NewTighteningIO(c *TighteningController) tightening_device.ITighteningIO {
	io := &TighteningIO{
		controller: c,
	}
	io.BaseDevice = device.CreateBaseDevice(device.BaseDeviceTypeIO, c.diag, nil, c.SerialNumber()) // 控制器io的序列号跟随控制器序列号
	return io
}

type TighteningIO struct {
	device.BaseDevice
	controller *TighteningController
}

func (s *TighteningIO) DeviceType() string {
	return device.BaseDeviceTypeIO
}

func (s *TighteningIO) SetIONotify(notify io.IONotify) {
	s.controller.SetIONotify(notify)
}

func (s *TighteningIO) IORead() (string, string, error) {
	return s.controller.IORead()
}

func (s *TighteningIO) IOWrite(index uint16, status uint16) error {
	return s.controller.IOWrite(index, status)
}
