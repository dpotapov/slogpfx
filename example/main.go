package main

import (
	"log/slog"
	"os"

	"github.com/dpotapov/slogpfx"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

func main() {
	h := tint.NewHandler(colorable.NewColorable(os.Stdout), nil)

	prefixed := slogpfx.NewHandler(h, &slogpfx.HandlerOptions{
		PrefixKeys: []string{"service", "geo"},
	})

	logger := slog.New(prefixed)

	logger.Info("Hello World!")
	logger.Info("Hello World!", "service", "billing")
	logger.With("service", "database").Error("Connection error")
	logger.With("service", "web", "geo", "eu").Warn("Access denied", "user", "admin", "action", "remove")
	logger.With("service", "web", "geo", "us").Info("Access granted", "user", "admin", "action", "create")
}
