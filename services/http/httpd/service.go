package httpd

import (
	stdContext "context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kevin0120/GoScrewdriverWebApi/services/diagnostic"
	"github.com/rs/cors"
	"log"
	"net/url"
	"strings"
	"time"
)

const (
	// BasePath Root path for the API
	BasePath = "/rush/v1"
)

var Httpd *Service

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
}

type Service struct {
	addr string
	err  chan error

	methods Methods

	ApiDoc          string
	Handler         []*Handler
	shutdownTimeout time.Duration
	externalURL     string
	server          *iris.Application

	stop chan chan struct{}

	HandlerByNames map[string]int

	cors   CorsConfig
	diag   Diagnostic
	opened bool

	DiagService interface {
		SetLogLevelFromName(lvl string) error
	}
	httpServerErrorLogger *log.Logger
}

func NewService(doc string, c Config, hostname string, d Diagnostic, disc *diagnostic.Service) (*Service, error) {

	port, _ := c.Port()
	u := url.URL{
		Host:   fmt.Sprintf("%s:%d", hostname, port),
		Scheme: "http",
	}
	s := &Service{
		ApiDoc:          doc,
		addr:            c.BindAddress,
		externalURL:     u.String(),
		cors:            c.Cors,
		server:          iris.New(),
		err:             make(chan error, 1),
		HandlerByNames:  make(map[string]int),
		shutdownTimeout: time.Duration(c.ShutdownTimeout),
		DiagService:     disc,
		opened:          false,
		diag:            d,
	}

	s.methods = newHttpMethods(s)

	if err := s.addNewHandler(BasePath, c, disc); err != nil {
		return nil, err
	}

	r := Route{
		RouteType:   ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/doc",
		HandlerFunc: s.methods.getDoc,
	}
	err := s.Handler[0].AddRoute(r)
	if err != nil {
		return nil, err
	}
	Httpd = s
	return s, nil
}

func (s *Service) addNewHandler(version string, c Config, disc *diagnostic.Service) error {
	if _, ok := s.HandlerByNames[version]; ok {
		// Should be unreachable code
		panic("cannot append handler twice")
	}
	crs := cors.New(cors.Options{
		AllowedOrigins:   s.cors.AllowedOrigins,
		AllowCredentials: s.cors.AllowCredentials,
		AllowedMethods:   s.cors.AllowedMethods,
	})
	s.server.WrapRouter(crs.ServeHTTP)
	p := s.server.Party(version).AllowMethods(iris.MethodOptions)
	if p == nil {
		return fmt.Errorf("fail to create the party %s", version)
	}
	h := NewHandler(
		c.LogEnabled,
		c.WriteTracing,
	)
	h.service = s.server
	h.DiagService = disc
	h.Version = version
	h.party = &p

	i := len(s.Handler)
	s.Handler = append(s.Handler, h)

	s.HandlerByNames[version] = i

	return nil
}

func (s *Service) manage() {
	//println("start mamager")
	var stopDone chan struct{}
	select {
	case stopDone = <-s.stop:
		// if we're already all empty, we're already done
		timeout := s.shutdownTimeout
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		s.server.Shutdown(ctx)
		close(stopDone)
		return
	}

}

// Close closes the underlying listener.
func (s *Service) Close() error {
	//defer s.DiagService.StoppedService()
	// If server is not set we were never started
	if s.server == nil {
		return nil
	}
	// Signal to manage loop we are stopping
	stopping := make(chan struct{})
	s.stop <- stopping

	<-stopping
	s.server = nil
	return nil
}

func (s *Service) serve() {
	s.diag.Info(s.URL())
	s.diag.Info(s.ExternalURL())
	err := s.server.Run(s.Addr(), iris.WithoutInterruptHandler)
	// The listener was closed so exit
	// See https://github.com/golang/go/issues/4373
	if !strings.Contains(err.Error(), "closed") {
		s.err <- fmt.Errorf("listener failed: addr=%s, err=%s", s.Addr(), err)
	} else {
		s.err <- nil
	}
}

// Open starts the service
func (s *Service) Open() error {
	//s.DiagService.StartingService()

	s.stop = make(chan chan struct{})

	go s.manage()
	go s.serve()
	return nil
}

func (s *Service) Addr() iris.Runner {
	return iris.Addr(s.addr)
}

func (s *Service) Err() <-chan error {
	return s.err
}

func (s *Service) URL() string {

	return "http://" + s.server.ConfigurationReadOnly().GetVHost()
}

// ExternalURL URL that should resolve externally to the server HTTP endpoint.
// It is possible that the URL does not resolve correctly  if the hostname config setting is incorrect.
func (s *Service) ExternalURL() string {
	return s.externalURL
}

func (s *Service) GetHandlerByName(version string) (*Handler, error) {
	i, ok := s.HandlerByNames[version]
	if !ok {
		// Should be unreachable code
		return nil, fmt.Errorf("cannot get handler By %s", version)
	}

	return s.Handler[i], nil
}

func (s *Service) AddNewHttpHandler(r Route) error {
	h, err := s.GetHandlerByName(BasePath)
	if err != nil {
		return err
	}
	err = h.AddRoute(r)
	if err != nil {
		return err
	}
	return nil
}
