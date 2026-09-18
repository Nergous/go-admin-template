package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

const (
	envFileName    = ".env"
	envFilePathEnv = "APP_ENV_FILE"
)

type Config struct {
	Environment Environment
	HTTP        HTTP
	DB          DB
	Logger      Logger
}

func Load() (*Config, error) {
	if err := loadDotEnv(); err != nil {
		return nil, fmt.Errorf("dotenv: %w", err)
	}

	var errs []error

	env, err := LoadEnvironment()
	if err != nil {
		errs = append(errs, fmt.Errorf("environment: %w", err))
	}

	http, err := LoadHTTP()
	if err != nil {
		errs = append(errs, fmt.Errorf("http: %w", err))
	}

	db, err := LoadDB()
	if err != nil {
		errs = append(errs, fmt.Errorf("db: %w", err))
	}

	logger, err := LoadLogger()
	if err != nil {
		errs = append(errs, fmt.Errorf("logger: %w", err))
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	cfg := &Config{
		Environment: env,
		HTTP:        http,
		DB:          db,
		Logger:      logger,
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var errs []error

	if c == nil {
		return errors.New("config is nil")
	}

	if err := c.Environment.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("environment: %w", err))
	}

	if err := c.HTTP.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("http: %w", err))
	}

	if err := c.DB.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("db: %w", err))
	}

	if err := c.Logger.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("logger: %w", err))
	}

	return errors.Join(errs...)
}

func loadDotEnv() error {
	explicitPath := strings.TrimSpace(os.Getenv(envFilePathEnv))
	if explicitPath != "" {
		return loadDotEnvFile(explicitPath, filepath.Abs, os.Stat, godotenv.Load)
	}

	if strings.EqualFold(strings.TrimSpace(os.Getenv(appEnvKey)), EnvironmentProd.String()) {
		return nil
	}

	candidates := []string{envFileName}
	if executable, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(executable), envFileName))
	}

	if err := loadDotEnvCandidates(candidates, filepath.Abs, os.Stat, godotenv.Load); err != nil {
		return err
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv(appEnvKey)), EnvironmentProd.String()) {
		return fmt.Errorf("implicit %s must not select production; set %s explicitly", envFileName, envFilePathEnv)
	}

	return nil
}

func loadDotEnvFile(
	candidate string,
	absolutePath func(string) (string, error),
	statFile func(string) (os.FileInfo, error),
	loadFile func(...string) error,
) error {
	absolute, err := absolutePath(candidate)
	if err != nil {
		return fmt.Errorf("resolve env file %q: %w", candidate, err)
	}

	info, err := statFile(absolute)
	if err != nil {
		return fmt.Errorf("inspect env file %q: %w", absolute, err)
	}
	if info.IsDir() {
		return fmt.Errorf("env file %q is a directory", absolute)
	}
	if err := loadFile(absolute); err != nil {
		return fmt.Errorf("load env file %q: %w", absolute, err)
	}

	return nil
}

func loadDotEnvCandidates(
	candidates []string,
	absolutePath func(string) (string, error),
	statFile func(string) (os.FileInfo, error),
	loadFile func(...string) error,
) error {
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		absolute, err := absolutePath(candidate)
		if err != nil {
			return fmt.Errorf("resolve env file %q: %w", candidate, err)
		}

		if _, ok := seen[absolute]; ok {
			continue
		}

		seen[absolute] = struct{}{}

		info, err := statFile(absolute)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return fmt.Errorf("inspect env file %q: %w", absolute, err)
		}

		if info.IsDir() {
			return fmt.Errorf("env file %q is a directory", absolute)
		}

		if err := loadFile(absolute); err != nil {
			return fmt.Errorf("load env file %q: %w", absolute, err)
		}

		return nil
	}

	return nil
}
