package tightening_device

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

const (
	ModelDesoutterCvi3        = "ModelDesoutterCvi3"
	ModelDesoutterCvi3Twin    = "ModelDesoutterCvi3Twin"
	ModelDesoutterCvi2        = "ModelDesoutterCvi2"
	ModelDesoutterDeltaWrench = "ModelDesoutterDeltaWrench"
	ModelDesoutterConnector   = "ModelDesoutterConnector"
	ModelCraneIQWrench        = "ModelCraneIQWrench"
	ModelLexenWrench          = "ModelLexenWrench"
	ModelLeetxTCS2000         = "ModelLeetxTCS2000"
)

type Service struct {
	configValue        atomic.Value
	runningControllers map[string]ITighteningController
	mtxDevices         sync.Mutex
	protocols          map[string]ITighteningProtocol
}

func (s *Service) loadTighteningController(c Config) {
	for k, deviceConfig := range c.Devices {
		p, err := s.getProtocol(deviceConfig.Protocol)
		if err != nil {
			continue
		}

		c, err := p.NewController(&deviceConfig) //如果不传index或导致获取的配置信息有误
		if err != nil {
			continue
		}

		sn := deviceConfig.SN
		if sn == "" {
			sn = fmt.Sprintf("%d", k+1)
		}

		//c.SetSerialNumber(sn)
		s.addController(sn, c)
	}
}

func NewService(c Config, protocols []ITighteningProtocol) (*Service, error) {

	s := &Service{
		runningControllers: map[string]ITighteningController{},
		protocols:          map[string]ITighteningProtocol{},
	}
	s.configValue.Store(c)
	// 载入支持的协议
	for _, protocol := range protocols {
		s.protocols[protocol.Name()] = protocol
	}

	// 根据配置加载所有拧紧控制器
	s.loadTighteningController(c)
	return s, nil
}

func (s *Service) getProtocol(protocolName string) (ITighteningProtocol, error) {
	if p, ok := s.protocols[protocolName]; !ok {
		return nil, errors.New("Protocol Is Not Support")
	} else {
		return p, nil
	}
}

//func (s *Service) setupHttpRoute() {
//
//	// @TODO 重复的接口
//	//r = httpd.Route{
//	//	RouteType:   httpd.ROUTE_TYPE_HTTP,
//	//	Method:      "PUT",
//	//	Pattern:     "/tool-enable",
//	//	HandlerFunc: s.putToolEnable,
//	//}
//	//s.httpd.AddNewHttpHandler(r)
//
//	r := httpd.Route{
//		RouteType:   httpd.ROUTE_TYPE_HTTP,
//		Method:      "PUT",
//		Pattern:     "/tool-pset",
//		HandlerFunc: s.putToolPSet,
//	}
//	if err := s.httpd.AddNewHttpHandler(r); err != nil {
//		s.diag.Error("AddNewHttpHandler tool-pset Error", err)
//	}
//}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	// 启动所有拧紧控制器
	s.startupControllers()

	return nil
}

func (s *Service) Close() error {

	// 关闭所有控制器
	s.shutdownControllers()

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}
func (s *Service) getControllers() map[string]ITighteningController {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	return s.runningControllers
}

func (s *Service) addController(controllerSN string, controller ITighteningController) {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	_, exist := s.runningControllers[controllerSN]
	if exist {
		return
	}

	s.runningControllers[controllerSN] = controller
}

func (s *Service) getController(controllerSN string) (ITighteningController, error) {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	td, exist := s.runningControllers[controllerSN]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Controller %s Not Found", controllerSN))
	}

	return td, nil
}

func (s *Service) getTool(controllerSN, toolSN string) (ITighteningTool, error) {
	_, err := s.getController(controllerSN)
	if err != nil {
		return nil, err
	}
	return nil, nil
	//tool, err := controller.GetToolViaSerialNumber(toolSN)
	//if err == nil {
	//	return tool, nil
	//}
	//tool, err = s.getFirstTool(controller)
	//
	//return tool, err
}

func (s *Service) getFirstTool(controller ITighteningController) (ITighteningTool, error) {
	if controller == nil {
		return nil, errors.New("getFirstTool: Controller Is Nil")
	}
	//
	//for _, v := range controller.Children() {
	//	return v.(ITighteningTool), nil
	//}

	return nil, errors.New("getFirstTool: Controller's Tool Not Found")
}

func (s *Service) startupControllers() {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	for _, c := range s.runningControllers {
		err := c.Start()
		if err != nil {
			continue
		}
	}
}

func (s *Service) shutdownControllers() {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	//for _, c := range s.runningControllers {
	//	err := c.Stop()
	//	if err != nil {
	//		continue
	//	}
	//}
}
