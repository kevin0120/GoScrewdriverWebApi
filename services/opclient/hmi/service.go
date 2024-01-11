package hmi

import (
	"github.com/kevin0120/GoScrewdriverWebApi/services/http/httpd"
	"go.uber.org/atomic"
)

type Service struct {
	diag        Diagnostic
	configValue atomic.Value

	httpd             HTTPService
	TighteningService ITightening

	workcenterSN string
}

func NewService(d Diagnostic, httpd HTTPService, sn string, tightening ITightening) *Service {

	s := &Service{
		diag:              d,
		httpd:             httpd,
		workcenterSN:      sn,
		TighteningService: tightening,
	}
	s.configValue.Store(Config{
		true,
	})

	s.setupTestInterface()

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) setupTestInterface() {
	var r httpd.Route
	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "POST",
		Pattern:     "/tighteningControl",
		HandlerFunc: s.tighteningControl,
	}
	s.httpd.AddNewHttpHandler(r)

}
