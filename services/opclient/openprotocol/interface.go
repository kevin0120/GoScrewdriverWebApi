package openprotocol

import (
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
)

type IClientHandler interface {
	SerialNumber() string
	//doDispatch(name string, data interface{})
	handleMsg(pkg *handlerPkg, context *clientContext) error
	HandleStatus(sn string, status string)
	GetVendorMid(mid string) (string, error)
	UpdateToolStatus(sn string, status string)
}

type IOpenProtocolController interface {
	tightening_device.ITighteningController

	//// 初始化控制器
	initController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service)
	InitSubscribeInfos()
	// GetMidHandler // Vendor Model定义(MID，IO等)
	GetVendorModel() map[string]interface{}
	//
	// op协议handler
	GetMidHandler(mid string) (MidHandler, error)
	//
	////控制器状态变化影响相关工具的状态变化
	//UpdateToolStatus(sn string, status string)
	//
	////处理未被处理的历史数据
	//handlerOldResults() error

	// 加载的协议
	Protocol() string

	// OpenProtocolParams New //曲线解析
	//CurveDataDecoding(original []byte, torqueCoefficient float64, angleCoefficient float64, d Diagnostic) (Torque []float32, Angle []float32)
	//
	OpenProtocolParams() *OpenProtocolParams
	// New
	//HandleStatus(sn string, status string)
	//
	New() IOpenProtocolController
}

type IResultData interface {
	ToTighteningResult() *tightening_device.TighteningResult
	GetInstance() interface{}
}
