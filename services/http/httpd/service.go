package httpd

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kevin0120/GoScrewdriverWebApi/services/diagnostic"
	"github.com/rs/cors"
	"log"
	"net/url"
	"time"
)

const (
	// BasePath Root path for the API
	BasePath = "/rush/v1"
)

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

	cors CorsConfig

	opened bool

	DiagService interface {
		SetLogLevelFromName(lvl string) error
	}
	httpServerErrorLogger *log.Logger
}

func NewService(doc string, c Config, hostname string, disc *diagnostic.Service) (*Service, error) {

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
	if err := s.AddNewHttpHandler(r); err != nil {
		return nil, err
	}

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
