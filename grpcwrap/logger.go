package grpcwrap

import (
	"fmt"
	"os"
	"runtime"

	"github.com/DataWorkbench/glog"
)

const (
	// InfoLevel indicates Info severity.
	InfoLevel int = iota + 1
	// WarningLevel indicates Warning severity.
	WarningLevel
	// ErrorLevel indicates Error severity.
	ErrorLevel
	// fatalLevel indicates Fatal severity.
	FatalLevel
)

// Logger for implements interface{} grpclog.LoggerV2
type Logger struct {
	Output    *glog.Logger
	Level     int
	Verbosity int
}

// NewLogger return a new Logger
func NewLogger(output *glog.Logger) *Logger {
	return &Logger{
		Output:    output,
		Level:     WarningLevel,
		Verbosity: 0,
	}
}

func (g *Logger) WithLevel(level int) *Logger {
	g.Level = level
	return g
}

func (g *Logger) WithVerbosity(v int) *Logger {
	g.Verbosity = v
	return g
}

func (g *Logger) WithOutput(l *glog.Logger) *Logger {
	g.Output = l
	return g
}

// implements grpclog.LoggerV2
//
// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (g *Logger) Info(args ...interface{}) {
	if InfoLevel < g.Level {
		return
	}
	g.Output.Info().RawString("grpclog-Info", fmt.Sprint(args...)).Fire()
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (g *Logger) Infoln(args ...interface{}) {
	if InfoLevel < g.Level {
		return
	}
	g.Output.Info().RawString("grpclog-Infoln", fmt.Sprint(args...)).Fire()
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (g *Logger) Infof(format string, args ...interface{}) {
	if InfoLevel < g.Level {
		return
	}
	g.Output.Info().RawString("grpclog-Infof", fmt.Sprintf(format, args...)).Fire()
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (g *Logger) Warning(args ...interface{}) {
	if WarningLevel < g.Level {
		return
	}
	g.Output.Warn().RawString("grpclog-Warning", fmt.Sprint(args...)).Fire()
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (g *Logger) Warningln(args ...interface{}) {
	if WarningLevel < g.Level {
		return
	}
	g.Output.Warn().RawString("grpclog-Warningln", fmt.Sprint(args...)).Fire()
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (g *Logger) Warningf(format string, args ...interface{}) {
	if WarningLevel < g.Level {
		return
	}
	g.Output.Warn().RawString("grpclog-Warningf", fmt.Sprintf(format, args...)).Fire()
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (g *Logger) Error(args ...interface{}) {
	if ErrorLevel < g.Level {
		return
	}
	g.Output.Error().RawString("grpclog-Error", fmt.Sprint(args...)).Fire()
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (g *Logger) Errorln(args ...interface{}) {
	if ErrorLevel < g.Level {
		return
	}
	g.Output.Error().RawString("grpclog-Errorln", fmt.Sprint(args...)).Fire()
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (g *Logger) Errorf(format string, args ...interface{}) {
	if ErrorLevel < g.Level {
		return
	}
	g.Output.Error().RawString("grpclog-Errorf", fmt.Sprintf(format, args...)).Fire()
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *Logger) Fatal(args ...interface{}) {
	if FatalLevel < g.Level {
		return
	}
	g.Output.Fatal().RawString("grpclog-Fatal", fmt.Sprint(args...)).Fire()
	os.Exit(1)
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *Logger) Fatalln(args ...interface{}) {
	if FatalLevel < g.Level {
		return
	}
	g.Output.Fatal().RawString("grpclog-Fatalln", fmt.Sprint(args...)).Fire()
	os.Exit(1)
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *Logger) Fatalf(format string, args ...interface{}) {
	if FatalLevel < g.Level {
		return
	}
	g.Output.Fatal().RawString("grpclog-Fatalf", fmt.Sprintf(format, args...)).Fire()
	os.Exit(1)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (g *Logger) V(l int) bool {
	return l <= g.Verbosity
}

// implements grpclog.DepthLoggerV2
//
// InfoDepth logs to INFO log at the specified depth. Arguments are handled in the manner of fmt.Print.
func (g *Logger) InfoDepth(depth int, args ...interface{}) {
	if InfoLevel < g.Level {
		return
	}
	_, file, line, ok := runtime.Caller(depth + 2)
	if ok {
		caller := fmt.Sprintf("%s:%d", file, line)
		g.Output.Info().RawString("grpclog-InfoDepth", fmt.Sprint(args...)).RawString("caller", caller).Fire()
	} else {
		g.Output.Info().RawString("grpclog-InfoDepth", fmt.Sprint(args...)).Fire()
	}
}

// WarningDepth logs to WARNING log at the specified depth. Arguments are handled in the manner of fmt.Print.
func (g *Logger) WarningDepth(depth int, args ...interface{}) {
	if WarningLevel < g.Level {
		return
	}
	_, file, line, ok := runtime.Caller(depth + 2)
	if ok {
		caller := fmt.Sprintf("%s:%d", file, line)
		g.Output.Info().RawString("grpclog-WarningDepth", fmt.Sprint(args...)).RawString("caller", caller).Fire()
	} else {
		g.Output.Warn().RawString("grpclog-WarningDepth", fmt.Sprint(args...)).Fire()
	}
}

// ErrorDetph logs to ERROR log at the specified depth. Arguments are handled in the manner of fmt.Print.
func (g *Logger) ErrorDepth(depth int, args ...interface{}) {
	if ErrorLevel < g.Level {
		return
	}
	_, file, line, ok := runtime.Caller(depth + 2)
	if ok {
		caller := fmt.Sprintf("%s:%d", file, line)
		g.Output.Info().RawString("grpclog-ErrorDepth", fmt.Sprint(args...)).RawString("caller", caller).Fire()
	} else {
		g.Output.Error().RawString("grpclog-ErrorDepth", fmt.Sprint(args...)).Fire()
	}
}

// FatalDepth logs to FATAL log at the specified depth. Arguments are handled in the manner of fmt.Print.
func (g *Logger) FatalDepth(depth int, args ...interface{}) {
	if FatalLevel < g.Level {
		return
	}
	_, file, line, ok := runtime.Caller(depth + 2)
	if ok {
		caller := fmt.Sprintf("%s:%d", file, line)
		g.Output.Info().RawString("grpclog-FatalDepth", fmt.Sprint(args...)).RawString("caller", caller).Fire()
	} else {
		g.Output.Fatal().RawString("grpclog-FatalDepth", fmt.Sprint(args...)).Fire()
	}
	os.Exit(1)
}
