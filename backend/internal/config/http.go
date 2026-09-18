package config

import (
	"errors"
	"fmt"
	"time"
)

const (
	httpPortEnv     = "HTTP_PORT"
	defaultHttpPort = 8080

	httpHostEnv     = "HTTP_HOST"
	defaultHttpHost = "127.0.0.1"

	httpTimeoutEnv = "HTTP_TIMEOUT"
	defaultTimeout = 30 * time.Second

	httpReadHeaderTimeoutEnv = "HTTP_READ_HEADER_TIMEOUT"
	defaultReadHeaderTimeout = 5 * time.Second

	httpReadTimeoutEnv  = "HTTP_READ_TIMEOUT"
	httpWriteTimeoutEnv = "HTTP_WRITE_TIMEOUT"
	httpIdleTimeoutEnv  = "HTTP_IDLE_TIMEOUT"
	defaultIdleTimeout  = 2 * time.Minute

	httpShutdownTimeoutEnv = "HTTP_SHUTDOWN_TIMEOUT"
	defaultShutdownTimeout = 10 * time.Second
)

type HTTP struct {
	Host              string
	Port              int
	Timeout           time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func LoadHTTP() (HTTP, error) {
	var errs []error

	host := stringEnv(defaultHttpHost, httpHostEnv)

	port, err := intEnv(defaultHttpPort, httpPortEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("port: %w", err))
	}

	timeout, err := timeEnv(defaultTimeout, httpTimeoutEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("timeout: %w", err))
	}

	readHeaderTimeout, err := timeEnv(defaultReadHeaderTimeout, httpReadHeaderTimeoutEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("readHeaderTimeout: %w", err))
	}

	readTimeout, err := timeEnv(timeout, httpReadTimeoutEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("readTimeout: %w", err))
	}

	writeTimeout, err := timeEnv(timeout, httpWriteTimeoutEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("writeTimeout: %w", err))
	}

	idleTimeout, err := timeEnv(defaultIdleTimeout, httpIdleTimeoutEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("idleTimeout: %w", err))
	}

	shutdownTimeout, err := timeEnv(defaultShutdownTimeout, httpShutdownTimeoutEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("shutdownTimeout: %w", err))
	}

	if err := errors.Join(errs...); err != nil {
		return HTTP{}, err
	}

	return NewHTTP(
		host,
		port,
		timeout,
		readHeaderTimeout,
		readTimeout,
		writeTimeout,
		idleTimeout,
		shutdownTimeout,
	)
}

func NewHTTP(
	host string,
	port int,
	timeout time.Duration,
	readHeaderTimeout time.Duration,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,
	shutdownTimeout time.Duration,
) (HTTP, error) {
	config := HTTP{
		Host:              normalizeHost(host),
		Port:              port,
		Timeout:           timeout,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ShutdownTimeout:   shutdownTimeout,
	}

	if err := config.Validate(); err != nil {
		return HTTP{}, err
	}

	return config, nil
}

func (c *HTTP) Validate() error {
	if c == nil {
		return errors.New("http is nil")
	}

	var errs []error

	if c.Host == "" {
		errs = append(errs, errors.New("host must not be empty"))
	} else if err := validateHost(c.Host); err != nil {
		errs = append(errs, fmt.Errorf("host %q: %w", c.Host, err))
	}

	if err := validatePort(c.Port); err != nil {
		errs = append(errs, err)
	}

	timeouts := []struct {
		name  string
		value time.Duration
	}{
		{name: "timeout", value: c.Timeout},
		{name: "readHeaderTimeout", value: c.ReadHeaderTimeout},
		{name: "readTimeout", value: c.ReadTimeout},
		{name: "writeTimeout", value: c.WriteTimeout},
		{name: "idleTimeout", value: c.IdleTimeout},
		{name: "shutdownTimeout", value: c.ShutdownTimeout},
	}
	for _, timeout := range timeouts {
		if timeout.value <= 0 {
			errs = append(errs, fmt.Errorf("%s must be greater than 0, got %v", timeout.name, timeout.value))
		}
	}

	return errors.Join(errs...)
}
