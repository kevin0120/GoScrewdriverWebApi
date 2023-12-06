package openprotocol

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
)

type IClientHandler interface {
	SerialNumber() string
	doDispatch(name string, data interface{})
	handleMsg(pkg *handlerPkg, context *clientContext) error
	HandleStatus(sn string, status string)
	GetVendorMid(mid string) (string, error)
	UpdateToolStatus(sn string, status string)
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
	ReleaseDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
}

type IStorageService interface {
	FindTargetResultForJobManual(workorderID int64) (storage.Results, error)
	UpdateTool(gun *storage.Tools) error
	ClearToolResultAndCurve(toolSN string) error
	GetTool(serial string) (storage.Tools, error)
	GetStep(id int64) (storage.Steps, error)
	UpdateIncompleteCurveAndSaveResult(result *storage.Results) error
	StorageInsertResult(result *storage.Results) error
	UpdateIncompleteResultAndSaveCurve(curve *storage.Curves) error
	PatchResultFromDB(result *storage.Results, mode string) error
	GetResultByTighteningID(toolSN string, tid string) (*storage.Results, error)
	UpdateRecord(bean interface{}, id int64, data map[string]interface{}) error
}

type IOpenProtocolController interface {
	tightening_device.ITighteningController

	// 初始化控制器
	initController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service, dp Dispatcher)
	InitSubscribeInfos()
	// Vendor Model定义(MID，IO等)
	GetVendorModel() map[string]interface{}

	// op协议handler
	GetMidHandler(mid string) (MidHandler, error)

	//控制器状态变化影响相关工具的状态变化
	UpdateToolStatus(sn string, status string)

	//处理未被处理的历史数据
	handlerOldResults() error

	// 加载的协议
	Protocol() string

	//曲线解析
	CurveDataDecoding(original []byte, torqueCoefficient float64, angleCoefficient float64, d Diagnostic) (Torque []float32, Angle []float32)

	OpenProtocolParams() *OpenProtocolParams

	HandleStatus(sn string, status string)

	New() IOpenProtocolController
}

type IResultData interface {
	ToTighteningResult() *tightening_device.TighteningResult
	GetInstance() interface{}
}
