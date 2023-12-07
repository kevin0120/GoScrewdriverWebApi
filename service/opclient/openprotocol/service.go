package openprotocol

import (
	"errors"
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/service/opclient/tightening_device"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	name        string
	vendors     map[string]IOpenProtocolController
}

func NewService(c Config, d Diagnostic, vendors map[string]IOpenProtocolController) *Service {

	s := &Service{
		name:    tightening_device.TIGHTENING_OPENPROTOCOL,
		diag:    d,
		vendors: vendors,
	}
	s.configValue.Store(c)
	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) NewController(cfg *tightening_device.TighteningDeviceConfig) (tightening_device.ITighteningController, error) {
	c, exist := s.vendors[cfg.Model]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Controller Model:%s Not Support", cfg.Model))
	}

	controllerInstance := c.New()
	controllerInstance.initController(cfg, s.diag, s)
	return controllerInstance, nil
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) GetDefaultMode() string {
	c := s.config()
	return c.DefaultMode
}

func (s *Service) generateIDInfo(info string) string {
	ids := ""
	for i := 0; i < MaxIdsNum; i++ {
		if i == s.config().DataIndex {
			ids += fmt.Sprintf("%-25s", info)
		} else {
			ids += fmt.Sprintf("%25s", "")
		}
	}

	return ids
}
