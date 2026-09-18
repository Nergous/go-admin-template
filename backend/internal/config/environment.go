package config

import (
	"fmt"
	"strings"
)

type Environment string

const (
	EnvironmentUnknown Environment = ""
	EnvironmentDev     Environment = "dev"
	EnvironmentProd    Environment = "prod"
	EnvironmentTest    Environment = "test"
)

const (
	appEnvKey          = "APP_ENV"
	defaultEnvironment = EnvironmentDev
)

func Environments() []Environment {
	return []Environment{
		EnvironmentDev,
		EnvironmentProd,
		EnvironmentTest,
	}
}

func LoadEnvironment() (Environment, error) {
	environment := stringEnv(defaultEnvironment.String(), appEnvKey)
	return NewEnvironment(environment)
}

func NewEnvironment(value string) (Environment, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	environment := Environment(value)
	if err := environment.Validate(); err != nil {
		return EnvironmentUnknown, err
	}

	return environment, nil
}

func (e Environment) String() string {
	return string(e)
}

func (e Environment) Validate() error {
	switch e {
	case EnvironmentDev, EnvironmentProd, EnvironmentTest:
		return nil
	default:
		return fmt.Errorf("invalid environment: %q", e)
	}
}
