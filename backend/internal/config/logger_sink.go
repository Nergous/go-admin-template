package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

type LogSink string

const (
	LogSinkUnknown LogSink = ""
	LogSinkStdout  LogSink = "stdout"
	LogSinkStderr  LogSink = "stderr"
)

func LogSinks() []LogSink {
	return []LogSink{
		LogSinkStdout,
		LogSinkStderr,
	}
}

const (
	logSinkEnv     = "LOG_SINK"
	defaultLogSink = LogSinkStdout
)

func LoadLogSink() (LogSink, error) {
	sink := stringEnv(defaultLogSink.String(), logSinkEnv)
	return NewLogSink(sink)
}

func NewLogSink(value string) (LogSink, error) {
	value = strings.TrimSpace(value)
	namedSink := LogSink(strings.ToLower(value))
	if namedSink == LogSinkStdout || namedSink == LogSinkStderr {
		return namedSink, nil
	}

	fileSink := LogSink(value)
	if err := fileSink.Validate(); err != nil {
		return defaultLogSink, err
	}

	return fileSink, nil
}

func (l LogSink) String() string {
	return string(l)
}

func (l LogSink) Validate() error {
	switch l {
	case LogSinkStdout, LogSinkStderr:
		return nil
	default:
		if !filepath.IsAbs(l.String()) {
			return fmt.Errorf("invalid log sink: %q", l)
		}
		return nil
	}
}
