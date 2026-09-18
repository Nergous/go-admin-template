package config

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"
)

func stringEnv(fallback, key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}

	return value
}

func rawEnv(fallback, key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}

func intEnv(fallback int, key string) (int, error) {
	value := stringEnv("", key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}

	return parsed, nil
}

func timeEnv(fallback time.Duration, key string) (time.Duration, error) {
	value := stringEnv("", key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a duration: %w", key, err)
	}

	return parsed, nil
}

func boolEnv(fallback bool, key string) (bool, error) {
	value := stringEnv("", key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}
	return parsed, nil
}

func validateHost(host string) error {
	if host == "" {
		return errors.New("empty host")
	}

	if _, err := netip.ParseAddr(normalizeHost(host)); err == nil {
		return nil
	}

	dnsHost := strings.TrimSuffix(host, ".")
	if len(dnsHost) > 253 {
		return fmt.Errorf("host too long: %d", len(host))
	}
	for label := range strings.SplitSeq(dnsHost, ".") {
		if label == "" || len(label) > 63 {
			return fmt.Errorf("invalid label %q in %q", label, host)
		}

		for i, r := range label {
			ok := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '-' || r == '_'
			if !ok {
				return fmt.Errorf("invalid char %q in %q", r, host)
			}

			if r == '-' && (i == 0 || i == len(label)-1) {
				return fmt.Errorf("label %q starts/ends with hyphen", label)
			}
		}
	}

	return nil
}

func normalizeHost(host string) string {
	host = strings.TrimSpace(host)
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return strings.TrimSuffix(strings.TrimPrefix(host, "["), "]")
	}

	return host
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port %d: must be between 1 and 65535", port)
	}

	return nil
}
