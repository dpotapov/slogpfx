package slogpfx

import (
	"context"
	"log/slog"
	"slices"
	"strings"
)

const (
	ansiCyan  = "\x1b[36m"
	ansiReset = "\x1b[0m"
)

type HandlerOptions struct {
	PrefixKeys []string // A list of keys used to fetch prefix values from the log record.

	// PrefixFormatter is a function to format the prefix of the log record.
	// If it's not set, the DefaultPrefixFormatter with ColorizePrefix wrapper is used.
	PrefixFormatter func(prefixes []slog.Value) string
}

// Handler is a custom slog handler that wraps another slog.Handler to prepend a prefix to the log
// messages. The prefix is sourced from the log record's attributes using the keys specified
// in PrefixKeys.
type Handler struct {
	Next slog.Handler // The next log handler in the chain.

	opts     HandlerOptions // Configuration options for this handler.
	prefixes []slog.Value   // Cached list of prefix values.
}

var _ slog.Handler = (*Handler)(nil)

// NewHandler creates a new prefix logging handler.
// The new handler will prepend a prefix sourced from the log record's attributes to each log
// message before passing the record to the next handler.
func NewHandler(next slog.Handler, opts *HandlerOptions) *Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}
	return &Handler{
		Next:     next,
		opts:     *opts,
		prefixes: make([]slog.Value, len(opts.PrefixKeys)),
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Next.Enabled(ctx, level)
}

// Handle processes a log record, prepending a prefix to its message if needed, and then passes the
// record to the next handler.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	prefixes := h.prefixes

	if r.NumAttrs() > 0 {
		nr := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
		attrs := make([]slog.Attr, 0, r.NumAttrs())
		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, a)
			return true
		})
		if p, changed := h.extractPrefixes(attrs); changed {
			nr.AddAttrs(attrs...)
			r = nr
			prefixes = p
		}
	}

	f := h.opts.PrefixFormatter
	if f == nil {
		f = ColorizePrefix(DefaultPrefixFormatter)
	}

	r.Message = f(prefixes) + r.Message

	return h.Next.Handle(ctx, r)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	p, _ := h.extractPrefixes(attrs)
	return &Handler{
		Next:     h.Next.WithAttrs(attrs),
		opts:     h.opts,
		prefixes: p,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		Next:     h.Next.WithGroup(name),
		opts:     h.opts,
		prefixes: h.prefixes,
	}
}

// extractPrefixes scans the attributes for keys specified in PrefixKeys.
// If found, their values are saved in a new prefix list.
// The original attribute list will be modified to remove the extracted prefix attributes.
func (h *Handler) extractPrefixes(attrs []slog.Attr) (prefixes []slog.Value, changed bool) {
	prefixes = h.prefixes
	for i, attr := range attrs {
		idx := slices.IndexFunc(h.opts.PrefixKeys, func(s string) bool { return s == attr.Key })
		if idx >= 0 {
			if !changed {
				// make a copy of prefixes:
				prefixes = make([]slog.Value, len(h.prefixes))
				copy(prefixes, h.prefixes)
			}
			prefixes[idx] = attr.Value
			attrs[i] = slog.Attr{} // remove the prefix attribute
			changed = true
		}
	}
	return
}

// DefaultPrefixFormatter constructs a prefix string by joining all detected prefix values using ":".
// A " > " suffix is added at the end of the prefix string.
func DefaultPrefixFormatter(prefixes []slog.Value) string {
	p := make([]string, 0, len(prefixes))
	for _, prefix := range prefixes {
		if prefix.Any() == nil || prefix.String() == "" {
			continue // skip empty prefixes
		}
		p = append(p, prefix.String())
	}
	if len(p) == 0 {
		return ""
	}
	return strings.Join(p, ":") + " > "
}

// ColorizePrefix wraps a prefix formatter function to colorize its output with cyan ANSI codes.
func ColorizePrefix(f func(prefixes []slog.Value) string) func(prefixes []slog.Value) string {
	return func(prefixes []slog.Value) string {
		p := f(prefixes)
		if p == "" {
			return ""
		}
		return ansiCyan + p + ansiReset
	}
}
