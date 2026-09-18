package config

import "log/slog"

const redacted = "[REDACTED]"

type Secret struct {
	value string
}

func NewSecret(value string) Secret {
	return Secret{value: value}
}

func (s Secret) Reveal() string {
	return s.value
}

func (Secret) String() string {
	return redacted
}

func (Secret) GoString() string {
	return redacted
}

func (Secret) LogValue() slog.Value {
	return slog.StringValue(redacted)
}
