package httpd

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/kevin0120/GoScrewdriverWebApi/utils/toml"
	"github.com/pkg/errors"
)

const (
	DefaultShutdownTimeout = toml.Duration(time.Second * 10)
)

type CorsConfig struct {
	AllowedOrigins   []string `yaml:"allowed-origins" json:"allowed-origins"`
	AllowCredentials bool     `yaml:"allow-credentials" json:"allow-credentials"`
	AllowedMethods   []string `yaml:"allowed-methods" json:"allowed-methods"`
}

type Config struct {
	BindAddress     string        `yaml:"bind-address"  json:"bind-address"`
	LogEnabled      bool          `yaml:"log-enabled" json:"log-enabled"`
	WriteTracing    bool          `yaml:"write-tracing" json:"write-tracing"`
	ShutdownTimeout toml.Duration `yaml:"shutdown-timeout" json:"shutdown-timeout"`
	Cors            CorsConfig    `yaml:"cors" json:"cors"`
	AccessLog       bool          `yaml:"access_log" json:"access_log"`
}

func NewConfig() Config {
	return Config{
		BindAddress:     ":8080",
		LogEnabled:      true,
		ShutdownTimeout: DefaultShutdownTimeout,
		Cors: CorsConfig{AllowedOrigins: []string{"*"},
			AllowCredentials: true,
			AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS"}},
		AccessLog: false,
	}
}

func (c Config) Validate() error {
	_, port, err := net.SplitHostPort(c.BindAddress)
	if err != nil {
		return errors.Wrapf(err, "invalid http bind address %s", c.BindAddress)
	}
	if port == "" {
		return errors.Wrapf(err, "invalid http bind address, no port specified %s", c.BindAddress)
	}
	if pn, err := strconv.ParseInt(port, 10, 64); err != nil {
		return errors.Wrapf(err, "invalid http bind address port %s", port)
	} else if pn > 65535 || pn < 0 {
		return fmt.Errorf("invalid http bind address port %d: out of range", pn)
	}

	return nil
}

// Determine HTTP port from BindAddress.
func (c Config) Port() (int, error) {
	if err := c.Validate(); err != nil {
		return -1, err
	}
	// Ignore errors since we already validated
	_, portStr, _ := net.SplitHostPort(c.BindAddress)
	port, _ := strconv.ParseInt(portStr, 10, 64)
	return int(port), nil
}
