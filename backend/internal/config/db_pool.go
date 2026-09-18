package config

import (
	"errors"
	"fmt"
	"time"
)

type DBPool struct {
	MaxOpen         int
	MaxIdle         int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func newDBPool(maxOpen int, maxIdle int, connMaxLifetime time.Duration, connMaxIdleTime time.Duration) (DBPool, error) {
	dbpool := DBPool{
		MaxOpen:         maxOpen,
		MaxIdle:         maxIdle,
		ConnMaxLifetime: connMaxLifetime,
		ConnMaxIdleTime: connMaxIdleTime,
	}

	if err := dbpool.Validate(); err != nil {
		return DBPool{}, err
	}

	return dbpool, nil
}

func (c *DBPool) Validate() error {
	if c == nil {
		return errors.New("dbpool is nil")
	}

	var errs []error

	if c.MaxOpen < 1 {
		errs = append(errs, fmt.Errorf("maxOpen must be greater than 0, got %v", c.MaxOpen))
	}
	if c.MaxIdle < 1 {
		errs = append(errs, fmt.Errorf("maxIdle must be greater than 0, got %v", c.MaxIdle))
	}
	if c.MaxOpen > 0 && c.MaxIdle > c.MaxOpen {
		errs = append(errs, fmt.Errorf("maxIdle must not exceed maxOpen, got %d > %d", c.MaxIdle, c.MaxOpen))
	}
	if c.ConnMaxLifetime < 0 {
		errs = append(errs, fmt.Errorf("connMaxLifetime must not be negative, got %v", c.ConnMaxLifetime))
	}
	if c.ConnMaxIdleTime < 0 {
		errs = append(errs, fmt.Errorf("connMaxIdleTime must not be negative, got %v", c.ConnMaxIdleTime))
	}

	return errors.Join(errs...)
}
