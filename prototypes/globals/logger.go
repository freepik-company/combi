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
	SLogger   *slog.Logger
	Context   context.Context
	ExtraAttr []any
}

var Logger LoggerT

func InitLogger(level LevelT, ctx context.Context, extraAttr ...any) {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.Level(level),
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	Logger.SLogger = slog.New(jsonHandler)
	Logger.Context = ctx
	Logger.ExtraAttr = extraAttr
}

func GetLevel(levelStr string) (l LevelT, err error) {
	levelMap := map[string]LevelT{
		"debug": DEBUG,
		"info":  INFO,
		"warn":  WARN,
		"error": ERROR,
	}
	if l, ok := levelMap[levelStr]; ok {
		return l, err
	}
	l = INFO
	err = fmt.Errorf("log level '%s' not supported", levelStr)
	return l, err
}

func (l *LoggerT) Debugf(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.DebugContext(l.Context, fmt.Sprintf(format, args...), l.ExtraAttr...)
		return
	}
	l.SLogger.Debug(fmt.Sprintf(format, args...), l.ExtraAttr...)
}

func (l *LoggerT) Infof(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.InfoContext(l.Context, fmt.Sprintf(format, args...), l.ExtraAttr...)
		return
	}
	l.SLogger.Info(fmt.Sprintf(format, args...), l.ExtraAttr...)
}

func (l *LoggerT) Warnf(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.WarnContext(l.Context, fmt.Sprintf(format, args...), l.ExtraAttr...)
		return
	}
	l.SLogger.Warn(fmt.Sprintf(format, args...), l.ExtraAttr...)
}

func (l *LoggerT) Errorf(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.ErrorContext(l.Context, fmt.Sprintf(format, args...), l.ExtraAttr...)
		return
	}
	l.SLogger.Error(fmt.Sprintf(format, args...), l.ExtraAttr...)
}

func (l *LoggerT) Fatalf(format string, args ...any) {
	if l.Context != nil {
		l.SLogger.ErrorContext(l.Context, fmt.Sprintf(format, args...), l.ExtraAttr...)
		os.Exit(1)
	}
	l.SLogger.Error(fmt.Sprintf(format, args...), l.ExtraAttr...)
	os.Exit(1)
}
