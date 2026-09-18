package main

import (
	"log/slog"
	"os"

	"go-admin-template/backend/internal/config"
	"go-admin-template/backend/internal/platform/observability"
)

func main() {
	bootstrapLogger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	cfg, err := config.Load()
	if err != nil {
		bootstrapLogger.Error("load configuration", "error", err)
		os.Exit(1)
	}

	logger, closeLogger, err := observability.NewWithConfig(cfg.Logger)
	if err != nil {
		bootstrapLogger.Error("create logger", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := closeLogger(); err != nil {
			bootstrapLogger.Error("close logger", "error", err)
		}
	}()

	logger.Info("application started")
}
