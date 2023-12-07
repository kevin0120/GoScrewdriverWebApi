package openprotocol

import "github.com/kevin0120/GoScrewdriverWebApi/service/opclient/tightening_device"

type ControllerSubscribe func(string) error

type handlerPkg struct {
	SN     string
	Header OpenProtocolHeader
	Body   string
	Seq    uint32
}

type respPkg struct {
	Seq  uint32
	Body interface{}
}

type SubscribeBarcodeStatusType string

type TighteningController struct {
	instance   IOpenProtocolController
	deviceConf *tightening_device.TighteningDeviceConfig
}

func (c *TighteningController) Protocol() string {
	return tightening_device.TIGHTENING_OPENPROTOCOL
}
func (c *TighteningController) model() string {
	return c.deviceConf.Model
	//return c.deviceConf.Model
}
func (c *TighteningController) Model() string {
	return c.model()
}

func (c *TighteningController) SetInstance(instance IOpenProtocolController) {
	c.instance = instance
}

func (c *TighteningController) initController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service) {
	c.deviceConf = deviceConfig
}
