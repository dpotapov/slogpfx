package slogpfx

import (
	"log/slog"
	"os"
)

func removeTimeAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{} // remove time attribute to avoid non-deterministic output
	}
	return a
}

func Example() {
	// Create a handler that writes to stdout.
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: removeTimeAttr})

	// Set the prefix for all log messages based on attribute "service".
	prefixed := NewHandler(h, &HandlerOptions{
		PrefixKeys:      []string{"service"},
		PrefixFormatter: DefaultPrefixFormatter,
	})

	logger := slog.New(prefixed)

	logger.Info("Hello World!")
	logger.Info("Hello World!", "service", "billing")
	logger.With("service", "database").Error("Connection error")

	// Output:
	// level=INFO msg="Hello World!"
	// level=INFO msg="billing > Hello World!"
	// level=ERROR msg="database > Connection error"
}

func Example_Multi() {
	// Create a handler that writes to stdout.
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: removeTimeAttr})

	// Set the prefix for all log messages based on attributes "service" and "geo".
	prefixed := NewHandler(h, &HandlerOptions{
		PrefixKeys:      []string{"service", "geo"},
		PrefixFormatter: DefaultPrefixFormatter,
	})

	logger := slog.New(prefixed)

	logger.Info("Hello World!")
	logger.Info("Hello World!", "service", "billing", "geo", "us")
	logger.With("service", "database", "geo", "eu").Error("Connection error")

	// Output:
	// level=INFO msg="Hello World!"
	// level=INFO msg="billing:us > Hello World!"
	// level=ERROR msg="database:eu > Connection error"
}
