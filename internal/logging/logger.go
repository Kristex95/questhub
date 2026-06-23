package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func New(level, format, logFile string) (*slog.Logger, error) {
	lvl, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{Level: lvl}

	switch strings.ToLower(format) {
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, opts)), nil
	case "text":
		return slog.New(slog.NewTextHandler(os.Stdout, opts)), nil
	case "ecs":
		handler, err := newECSHandler(logFile, opts)
		if err != nil {
			return nil, err
		}
		return slog.New(handler), nil
	default:
		return nil, fmt.Errorf("unknown log format %q (want json, text, or ecs)", format)
	}
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level %q", level)
	}
}

func newECSHandler(logFile string, opts *slog.HandlerOptions) (slog.Handler, error) {
	if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	return &ecsHandler{
		w:    io.MultiWriter(os.Stdout, file),
		opts: opts,
	}, nil
}

type ecsHandler struct {
	w    io.Writer
	opts *slog.HandlerOptions
	mu   sync.Mutex
}

func (h *ecsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.opts != nil && h.opts.Level != nil {
		return level >= h.opts.Level.Level()
	}
	return level >= slog.LevelInfo
}

func (h *ecsHandler) Handle(ctx context.Context, r slog.Record) error {
	entry := map[string]any{
		"@timestamp": r.Time.UTC().Format(time.RFC3339Nano),
		"log.level":  r.Level.String(),
		"message":    r.Message,
	}

	r.Attrs(func(a slog.Attr) bool {
		entry[a.Key] = a.Value.Any()
		return true
	})

	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		entry["trace.id"] = span.SpanContext().TraceID().String()
		entry["span.id"] = span.SpanContext().SpanID().String()
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err = h.w.Write(append(data, '\n'))
	return err
}

func (h *ecsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ecsHandler) WithGroup(_ string) slog.Handler {
	return h
}
