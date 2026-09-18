package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

const (
	dbDriverEnv = "DB_DRIVER"
	dbURLEnv    = "DB_URL"

	dbHostEnv     = "DB_HOST"
	dbPortEnv     = "DB_PORT"
	dbDatabaseEnv = "DB_DATABASE"
	dbUsernameEnv = "DB_USERNAME"
	dbPasswordEnv = "DB_PASSWORD"

	dbTimeoutEnv     = "DB_TIMEOUT"
	defaultDBTimeout = 10 * time.Second

	dbMaxOpenEnv     = "DB_MAX_OPEN"
	defaultDBMaxOpen = 10

	dbMaxIdleEnv     = "DB_MAX_IDLE"
	defaultDBMaxIdle = 10

	dbConnMaxLifetimeEnv = "DB_CONN_MAX_LIFETIME"
	dbConnMaxIdleTimeEnv = "DB_CONN_MAX_IDLE_TIME"
)

type DBDriver string

type DB struct {
	Driver     DBDriver
	Connection DBConnection
	Timeout    time.Duration
	Pool       DBPool
}

func LoadDB() (DB, error) {
	var errs []error

	driver := DBDriver(stringEnv("", dbDriverEnv))
	url := rawEnv("", dbURLEnv)
	host := stringEnv("", dbHostEnv)

	port := 0
	var err error
	if url == "" {
		port, err = intEnv(0, dbPortEnv)
		if err != nil {
			errs = append(errs, fmt.Errorf("port: %w", err))
		}
	}

	database := stringEnv("", dbDatabaseEnv)
	username := stringEnv("", dbUsernameEnv)
	password := rawEnv("", dbPasswordEnv)

	timeout, err := timeEnv(defaultDBTimeout, dbTimeoutEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("timeout: %w", err))
	}

	maxOpen, err := intEnv(defaultDBMaxOpen, dbMaxOpenEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("maxOpen: %w", err))
	}

	maxIdle, err := intEnv(defaultDBMaxIdle, dbMaxIdleEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("maxIdle: %w", err))
	}

	connMaxLifetime, err := timeEnv(0, dbConnMaxLifetimeEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("connMaxLifetime: %w", err))
	}

	connMaxIdleTime, err := timeEnv(0, dbConnMaxIdleTimeEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("connMaxIdleTime: %w", err))
	}

	if err := errors.Join(errs...); err != nil {
		return DB{}, err
	}

	dbconn, err := newDBConnection(url, host, port, database, username, password)
	if err != nil {
		return DB{}, err
	}

	dbpool, err := newDBPool(maxOpen, maxIdle, connMaxLifetime, connMaxIdleTime)
	if err != nil {
		return DB{}, err
	}

	return NewDB(driver, dbconn, timeout, dbpool)
}

func NewDB(driver DBDriver, connection DBConnection, timeout time.Duration, pool DBPool) (DB, error) {
	config := DB{
		Driver:     DBDriver(strings.ToLower(strings.TrimSpace(string(driver)))),
		Connection: connection,
		Timeout:    timeout,
		Pool:       pool,
	}

	if err := config.Validate(); err != nil {
		return DB{}, err
	}

	return config, nil
}

func (c *DB) String() string {
	return fmt.Sprintf(
		"DB{Driver:%q URL:%s Host:%q Port:%d Database:%q Username:%q Password:%s Timeout:%s MaxOpen:%d MaxIdle:%d ConnMaxLifetime:%s ConnMaxIdleTime:%s}",
		c.Driver, c.Connection.URL.String(), c.Connection.Host, c.Connection.Port, c.Connection.Database, c.Connection.Username, c.Connection.Password.String(),
		c.Timeout, c.Pool.MaxOpen, c.Pool.MaxIdle, c.Pool.ConnMaxLifetime, c.Pool.ConnMaxIdleTime,
	)
}

func (c *DB) GoString() string {
	return c.String()
}

func (c *DB) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("driver", string(c.Driver)),
		slog.String("url", c.Connection.URL.String()),
		slog.String("host", c.Connection.Host),
		slog.Int("port", c.Connection.Port),
		slog.String("database", c.Connection.Database),
		slog.String("username", c.Connection.Username),
		slog.String("password", c.Connection.Password.String()),
		slog.Duration("timeout", c.Timeout),
		slog.Int("max_open", c.Pool.MaxOpen),
		slog.Int("max_idle", c.Pool.MaxIdle),
		slog.Duration("conn_max_lifetime", c.Pool.ConnMaxLifetime),
		slog.Duration("conn_max_idle_time", c.Pool.ConnMaxIdleTime),
	)
}

func (c *DB) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}

	var errs []error

	if err := validateDBDriver(c.Driver); err != nil {
		errs = append(errs, err)
	}
	if err := c.Connection.Validate(); err != nil {
		errs = append(errs, err)
	}
	if c.Timeout < 1 {
		errs = append(errs, fmt.Errorf("timeout must be greater than 0, got %v", c.Timeout))
	}
	if err := c.Pool.Validate(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func validateDBDriver(driver DBDriver) error {
	if driver == "" {
		return errors.New("driver must not be empty")
	}

	for i, char := range driver {
		valid := char >= 'a' && char <= 'z' || char >= '0' && char <= '9' || char == '_' || char == '-'
		if !valid || i == 0 && char >= '0' && char <= '9' {
			return fmt.Errorf("invalid driver %q", driver)
		}
	}

	return nil
}
