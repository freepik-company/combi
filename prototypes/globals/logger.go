package globals

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// ----------------------------------------------------------------
// LOGGER
// ----------------------------------------------------------------

type LevelT int

const (
	DEBUG LevelT = LevelT(slog.LevelDebug)
	INFO  LevelT = LevelT(slog.LevelInfo)
	WARN  LevelT = LevelT(slog.LevelWarn)
	ERROR LevelT = LevelT(slog.LevelError)
)

type LoggerT struct {
	SLogger *slog.Logger
	Context context.Context
}

var Logger LoggerT

func InitLogger(level LevelT, ctx context.Context) {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.Level(level),
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	Logger.Context = ctx
	Logger.SLogger = slog.New(jsonHandler)
}

func (l *LoggerT) Debugf(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.DebugContext(l.Context, fmt.Sprintf(format, args...))
		return
	}
	l.SLogger.Debug(fmt.Sprintf(format, args...))
}

func (l *LoggerT) Infof(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.InfoContext(l.Context, fmt.Sprintf(format, args...))
		return
	}
	l.SLogger.Info(fmt.Sprintf(format, args...))
}

func (l *LoggerT) Warnf(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.WarnContext(l.Context, fmt.Sprintf(format, args...))
		return
	}
	l.SLogger.Warn(fmt.Sprintf(format, args...))
}

func (l *LoggerT) Errorf(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.ErrorContext(l.Context, fmt.Sprintf(format, args...))
		return
	}
	l.SLogger.Error(fmt.Sprintf(format, args...))
}

func (l *LoggerT) Fatalf(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.ErrorContext(l.Context, fmt.Sprintf(format, args...))
		os.Exit(1)
	}
	l.SLogger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
