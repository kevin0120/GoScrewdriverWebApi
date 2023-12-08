package openprotocol

import (
	"errors"
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/service/opclient/tightening_device"
	"go.uber.org/atomic"
	"time"
)

const (
	BaseDeviceStatusOnline  = "online"
	BaseDeviceStatusOffline = "offline"
	//拧紧工具特有的
	BaseDeviceStatusRunning   = "running"
	BaseDeviceStatusUnRunning = "stopping"
	BaseDeviceStatusEnabled   = "enabled"
	BaseDeviceStatusDisabled  = "disabled"
)

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
	instance             IOpenProtocolController
	deviceConf           *tightening_device.TighteningDeviceConfig
	sockClients          map[string]*clientContext
	ControllerSubscribes []ControllerSubscribe
	isGlobalConn         bool
	opened               atomic.Bool
	status               *atomic.String
	serialNumber         string
}

func (s *TighteningController) Protocol() string {
	return tightening_device.TIGHTENING_OPENPROTOCOL
}
func (s *TighteningController) model() string {
	return s.deviceConf.Model
	//return c.deviceConf.Model
}
func (s *TighteningController) Model() string {
	return s.model()
}

func (s *TighteningController) SetInstance(instance IOpenProtocolController) {
	s.instance = instance
}
func (s *TighteningController) getInstance() IOpenProtocolController {
	if s.instance == nil {
		panic("Controller Instance Is Nil")
	}

	return s.instance
}
func (s *TighteningController) New() IOpenProtocolController {
	//fixme: 永远不能被调用
	return nil
}
func (s *TighteningController) ResultSubscribe(sn string) error {
	//FIXME: 现在临时通过异常捕捉的方式进行修复
	defer func() { //进行异常捕捉
		if err := recover(); err != nil {
			fmt.Printf("ResultSubscribe error: %v", err)
		}
	}()

	//reply, err := c.getClient(sn).ProcessRequest(MID_0060_LAST_RESULT_SUBSCRIBE, false, "", "", "")
	//if err != nil {
	//	return err
	//}
	//
	//tt := reflect.TypeOf(reply)
	//ss := ""
	//switch tt.Kind() { //nolint: exhaustive
	//case reflect.String:
	//	ss = reply.(string)
	//default:
	//	c.diag.Error("ResultSubscribe", errors.New("reply Type Is Not String"))
	//}
	//
	//if ss != requestErrors["00"] {
	//	return errors.New(fmt.Sprintf("MID: %s err: %s", MID_0060_LAST_RESULT_SUBSCRIBE, reply.(string)))
	//}

	return nil
}
func (s *TighteningController) InitSubscribeInfos() {
	s.ControllerSubscribes = []ControllerSubscribe{
		s.ResultSubscribe,
		////c.SelectorSubscribe,
		//c.JobInfoSubscribe,
		//c.IOInputSubscribe,
		//c.VinSubscribe,
		//c.AlarmSubscribe,
		//c.CurveSubscribe,
	}
}

func (s *TighteningController) initController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service) {
	s.deviceConf = deviceConfig
	s.status = atomic.NewString(BaseDeviceStatusOffline)
	s.sockClients = map[string]*clientContext{}
	s.getInstance().InitSubscribeInfos()
	s.initClients(deviceConfig, d)
}

func (s *TighteningController) OpenProtocolParams() *OpenProtocolParams {
	return &OpenProtocolParams{
		MaxKeepAliveCheck: 3,
		MaxReplyTime:      3 * time.Second,
		KeepAlivePeriod:   8 * time.Second,
	}
}
func (s *TighteningController) initClients(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic) {

	for _, v := range deviceConfig.Tools {
		endpoint := v.Endpoint
		sn := v.SN
		if deviceConfig.Endpoint != "" {
			// 全局链接
			s.isGlobalConn = true
			endpoint = deviceConfig.Endpoint
			sn = deviceConfig.SN
		} else {
			// 每个工具独立链接
			s.isGlobalConn = false
		}

		client := newClientContext(endpoint, d, s.getInstance().(IClientHandler), sn, s.getInstance().OpenProtocolParams())
		s.sockClients[sn] = client

		if s.isGlobalConn {
			break
		}
	}
}

func (c *TighteningController) UpdateToolStatus(sn string, status string) {
	//tool, err := c.getToolViaSerialNumber(sn)
	//if err != nil {
	//	tool, _ = c.getInstance().GetToolViaChannel(1)
	//	//return
	//}
	//if tool == nil {
	//	return
	//}
	//tool.(*TighteningTool).UpdateStatus(status)
}
func (s *TighteningController) Status() string {
	return s.status.Load()
}

func (s *TighteningController) UpdateStatus(status string) {
	s.status.Store(status)
}

func (s *TighteningController) SerialNumber() string {
	return s.serialNumber
}

func (s *TighteningController) SetSerialNumber(serialNumber string) {
	s.serialNumber = serialNumber
}

func (s *TighteningController) HandleStatus(sn string, status string) {
	if status == s.Status() {
		return
	}
	s.UpdateStatus(status)
}
func (s *TighteningController) handleMsg(pkg *handlerPkg, context *clientContext) error {

	handler, err := s.getInstance().GetMidHandler(pkg.Header.MID)
	if err != nil {
		return err
	}

	return handler(s, pkg)
}

func (s *TighteningController) Start() error {
	if s.opened.Load() {
		return nil
	}
	s.opened.Store(true)

	// 启动客户端
	s.startupClients()

	return nil
}
func (s *TighteningController) startupClients() {

	for _, v := range s.sockClients {
		v.start()
	}
}
func (s *TighteningController) GetVendorMid(mid string) (string, error) {
	rev, exist := s.getInstance().GetVendorModel()[mid]
	if !exist {
		return "", errors.New(fmt.Sprintf("MID %s Not Supported", mid))
	}

	return rev.(string), nil
}
