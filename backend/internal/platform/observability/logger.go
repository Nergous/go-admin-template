package observability

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"go-admin-template/backend/internal/config"
)

func New() *slog.Logger {
	return slog.New(contextHandler{
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	})
}

func NewWithConfig(cfg config.Logger) (*slog.Logger, func() error, error) {
	if err := cfg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate logger config: %w", err)
	}

	options := &slog.HandlerOptions{
		Level:     cfg.Level.ToSlogLevel(),
		AddSource: cfg.AddSource,
	}

	output, closeOutput, err := sinkWriter(cfg.Sink)
	if err != nil {
		return nil, nil, err
	}

	handler := cfg.Format.SlogHandlerFunc()(output, options)

	return slog.New(contextHandler{Handler: handler}), closeOutput, nil
}

func sinkWriter(sink config.LogSink) (io.Writer, func() error, error) {
	noop := func() error {
		return nil
	}

	switch sink {
	case config.LogSinkUnknown, config.LogSinkStdout:
		return os.Stdout, noop, nil
	case config.LogSinkStderr:
		return os.Stderr, noop, nil
	}

	if err := sink.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate log sink: %w", err)
	}

	path := sink.String()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, nil, fmt.Errorf("create log directory %q: %w", filepath.Dir(path), err)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file %q: %w", path, err)
	}

	return file, file.Close, nil
}
