package hmi

import (
	"github.com/kataras/iris/v12"
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

func NewService(c Config, d Diagnostic, httpd HTTPService, sn string, tightening ITightening) *Service {

	s := &Service{
		diag:              d,
		httpd:             httpd,
		workcenterSN:      sn,
		TighteningService: tightening,
	}
	s.configValue.Store(c)

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
	s.diag.Close()
	s.diag.Closed()

	return nil
}

// healthz
func (s *Service) getHealthz(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNoContent)
}

func (s *Service) setupTestInterface() {
	var r httpd.Route

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/healthz",
		HandlerFunc: s.getHealthz,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/notify",
		HandlerFunc: s.putNotify,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/workorders",
		HandlerFunc: s.listWorkorders,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/workorder",
		HandlerFunc: s.getWorkorderDetail,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/local-results",
		HandlerFunc: s.getLocalResults,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/next-workorder",
		HandlerFunc: s.getNextWorkorder,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/tool-enable",
		HandlerFunc: s.putToolControl,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/psets",
		HandlerFunc: s.putPSets,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "POST",
		Pattern:     "/ak2",
		HandlerFunc: s.postAK2,
	}
	s.httpd.AddNewHttpHandler(r)

}
