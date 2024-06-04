package globals

import (
	"log/slog"
	"os"
)

// ----------------------------------------------------------------
// LOGGER
// ----------------------------------------------------------------

const (
	DEBUG = slog.LevelDebug
	INFO  = slog.LevelInfo
	WARN  = slog.LevelWarn
	ERROR = slog.LevelError
)

var Logger *slog.Logger

func InitLogger(level slog.Level) {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	Logger = slog.New(jsonHandler)
}
