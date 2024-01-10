package hmi

import (
	"github.com/kevin0120/GoScrewdriverWebApi/services/http/httpd"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
	Disconnect(id string)
	Close()
	Closed()
}

type HTTPService interface {
	AddNewHttpHandler(r httpd.Route) error
}

type ITightening interface {
	ToolPSetByIP(req *tightening_device.PSetSet) error
	ToolLedControl(toolSN string, enable bool) error
	ToolControl(req *tightening_device.ToolControl) error
	ToolPSetBatchSet(req *tightening_device.PSetBatchSet) error
	ToolPSetSet(req *tightening_device.PSetSet) error
	GetControllerByToolSN(toolSN string) (tightening_device.ITighteningController, error)
	GetToolByToolSN(toolSN string) (tightening_device.ITighteningTool, error)
}
