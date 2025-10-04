package log

import (
	"errors"
	"io"
	"log/slog"

	"github.com/lmittmann/tint"
)

func New(w io.Writer, version string) (*slog.Logger, *slog.LevelVar) {
	level := &slog.LevelVar{}
	return slog.New(tint.NewHandler(w, &tint.Options{
		Level: level,
	})).With("program", "ghir", "version", version), level
}

var ErrUnknownLogLevel = errors.New("unknown log level")

func SetLevel(levelVar *slog.LevelVar, level string) error {
	lvl, err := parseLevel(level)
	if err != nil {
		return err
	}
	levelVar.Set(lvl)
	return nil
}

func parseLevel(lvl string) (slog.Level, error) {
	switch lvl {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, ErrUnknownLogLevel
	}
}
