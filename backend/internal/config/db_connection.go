package config

import (
	"errors"
	"fmt"
	"strings"
)

type DBConnection struct {
	URL      Secret
	Host     string
	Port     int
	Database string
	Username string
	Password Secret
}

func newDBConnection(url, host string, port int, database, username, password string) (DBConnection, error) {

	dbconn := DBConnection{
		URL:      NewSecret(url),
		Host:     strings.TrimSpace(host),
		Port:     port,
		Database: strings.TrimSpace(database),
		Username: strings.TrimSpace(username),
		Password: NewSecret(password),
	}

	if err := dbconn.Validate(); err != nil {
		return DBConnection{}, err
	}

	return dbconn, nil
}

func (c *DBConnection) UsesURL() bool {
	return c.URL.Reveal() != ""
}

func (c *DBConnection) Validate() error {
	if c == nil {
		return errors.New("dbconnection is nil")
	}

	if c.UsesURL() {
		return nil
	}

	var errs []error

	if c.Database == "" {
		errs = append(errs, errors.New("DB_URL or DB_DATABASE must be set"))
	}
	if c.Port < 0 || c.Port > 65535 {
		errs = append(errs, fmt.Errorf("port %d: must be between 0 and 65535", c.Port))
	}

	return errors.Join(errs...)
}
