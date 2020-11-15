package gormwrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DataWorkbench/glog"
	"gorm.io/gorm/logger"
)

// LogLevel
type LogLevel = logger.LogLevel

const (
	SilentLevel LogLevel = iota + 1
	ErrorLevel
	WarnLevel
	InfoLevel
)

// Logger implements gorm logger.Interface used by glog.Logger
type Logger struct {
	Level         LogLevel
	SlowThreshold time.Duration
	Output        *glog.Logger // the default logger instances
}

func NewLogger(output *glog.Logger) *Logger {
	return &Logger{
		Level:         WarnLevel,
		SlowThreshold: time.Second * 2,
		Output:        output,
	}
}

func (g *Logger) WithLevel(level LogLevel) *Logger {
	g.Level = level
	return g
}

func (g *Logger) WithSlowThreshold(d time.Duration) *Logger {
	g.SlowThreshold = d
	return g
}

func (g *Logger) WithOutput(l *glog.Logger) *Logger {
	g.Output = l
	return g
}

func (g *Logger) LogMode(level LogLevel) logger.Interface {
	nl := *g
	nl.Level = level
	return &nl
}

func (g *Logger) Info(ctx context.Context, format string, v ...interface{}) {
	if g.Level < InfoLevel {
		return
	}
	l := glog.FromContext(ctx)
	if l == nil {
		l = g.Output
	}
	format = strings.TrimSuffix(format, "\n")
	l.Info().RawString("gormlog", fmt.Sprintf(format, v...)).Fire()
}

func (g *Logger) Warn(ctx context.Context, format string, v ...interface{}) {
	if g.Level < WarnLevel {
		return
	}
	l := glog.FromContext(ctx)
	if l == nil {
		l = g.Output
	}
	format = strings.TrimSuffix(format, "\n")
	l.Warn().RawString("gormlog", fmt.Sprintf(format, v...)).Fire()
}

func (g *Logger) Error(ctx context.Context, format string, v ...interface{}) {
	if g.Level < ErrorLevel {
		return
	}
	l := glog.FromContext(ctx)
	if l == nil {
		l = g.Output
	}
	format = strings.TrimSuffix(format, "\n")
	l.Error().RawString("gormlog", fmt.Sprintf(format, v...)).Fire()
}

func (g *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.Level < SilentLevel {
		return
	}
	l := glog.FromContext(ctx)
	if l == nil {
		l = g.Output
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && g.Level >= ErrorLevel:
		sql, rows := fc()
		l.Error().Msg("gormlog trace ").String("SQL", sql).Int64("rows", rows).Error("error", err).Millisecond("elapsed", elapsed).Fire()
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.Level >= WarnLevel:
		sql, rows := fc()
		l.Warn().Msg("gormlog trace ").String("SQL", sql).Int64("rows", rows).Millisecond("elapsed", elapsed).Fire()
	case g.Level >= InfoLevel:
		sql, rows := fc()
		l.Debug().Msg("gormlog trace ").String("SQL", sql).Int64("rows", rows).Millisecond("elapsed", elapsed).Fire()
	}
}
