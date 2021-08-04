package grpcwrap

import (
	"fmt"
	"os"
	"runtime"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc/grpclog"
)

const (
	// infoLevel indicates Info severity.
	infoLevel int = iota + 1
	// warningLevel indicates Warning severity.
	warningLevel
	// errorLevel indicates Error severity.
	errorLevel
	// fatalLevel indicates Fatal severity.
	fatalLevel
)

// SetLogger sets logger object that used in grpc.
// Not mutex-protected, should be called before any gRPC functions.
func SetLogger(output *glog.Logger, cfg *LogConfig) {
	grpclog.SetLoggerV2(&loggerT{
		output: output,
		cfg:    cfg,
	})
}

type LogConfig struct {
	// grpc log level: 1 => info, 2 => waring, 3 => error, 4 => fatal
	Level     int `json:"level"     yaml:"level"     env:"LEVEL,default=3"     validate:"gte=1,lte=4"`
	Verbosity int `json:"verbosity" yaml:"verbosity" env:"VERBOSITY,default=99" validate:"required"`
}

// loggerT for implements interface{} grpclog.LoggerV2 and grpclog.DepthLoggerV2
type loggerT struct {
	output *glog.Logger
	cfg    *LogConfig
}

// implements grpclog.LoggerV2
//
// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (g *loggerT) Info(args ...interface{}) {
	if infoLevel < g.cfg.Level {
		return
	}
	g.output.Info().RawString("grpclog-Info", fmt.Sprint(args...)).Fire()
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (g *loggerT) Infoln(args ...interface{}) {
	if infoLevel < g.cfg.Level {
		return
	}
	g.output.Info().RawString("grpclog-Infoln", fmt.Sprint(args...)).Fire()
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (g *loggerT) Infof(format string, args ...interface{}) {
	if infoLevel < g.cfg.Level {
		return
	}
	g.output.Info().RawString("grpclog-Infof", fmt.Sprintf(format, args...)).Fire()
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (g *loggerT) Warning(args ...interface{}) {
	if warningLevel < g.cfg.Level {
		return
	}
	g.output.Warn().RawString("grpclog-Warning", fmt.Sprint(args...)).Fire()
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (g *loggerT) Warningln(args ...interface{}) {
	if warningLevel < g.cfg.Level {
		return
	}
	g.output.Warn().RawString("grpclog-Warningln", fmt.Sprint(args...)).Fire()
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (g *loggerT) Warningf(format string, args ...interface{}) {
	if warningLevel < g.cfg.Level {
		return
	}
	g.output.Warn().RawString("grpclog-Warningf", fmt.Sprintf(format, args...)).Fire()
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (g *loggerT) Error(args ...interface{}) {
	if errorLevel < g.cfg.Level {
		return
	}
	g.output.Error().RawString("grpclog-Error", fmt.Sprint(args...)).Fire()
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (g *loggerT) Errorln(args ...interface{}) {
	if errorLevel < g.cfg.Level {
		return
	}
	g.output.Error().RawString("grpclog-Errorln", fmt.Sprint(args...)).Fire()
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (g *loggerT) Errorf(format string, args ...interface{}) {
	if errorLevel < g.cfg.Level {
		return
	}
	g.output.Error().RawString("grpclog-Errorf", fmt.Sprintf(format, args...)).Fire()
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *loggerT) Fatal(args ...interface{}) {
	if fatalLevel < g.cfg.Level {
		return
	}
	g.output.Fatal().RawString("grpclog-Fatal", fmt.Sprint(args...)).Fire()
	os.Exit(1)
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *loggerT) Fatalln(args ...interface{}) {
	if fatalLevel < g.cfg.Level {
		return
	}
	g.output.Fatal().RawString("grpclog-Fatalln", fmt.Sprint(args...)).Fire()
	os.Exit(1)
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *loggerT) Fatalf(format string, args ...interface{}) {
	if fatalLevel < g.cfg.Level {
		return
	}
	g.output.Fatal().RawString("grpclog-Fatalf", fmt.Sprintf(format, args...)).Fire()
	os.Exit(1)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (g *loggerT) V(l int) bool {
	return l <= g.cfg.Verbosity
}

// implements grpclog.DepthLoggerV2
//
// InfoDepth logs to INFO log at the specified depth. Arguments are handled in the manner of fmt.Print.
func (g *loggerT) InfoDepth(depth int, args ...interface{}) {
	if infoLevel < g.cfg.Level {
		return
	}
	_, file, line, ok := runtime.Caller(depth + 2)
	if ok {
		caller := fmt.Sprintf("%s:%d", file, line)
		g.output.Info().RawString("grpclog-InfoDepth", fmt.Sprint(args...)).RawString("caller", caller).Fire()
	} else {
		g.output.Info().RawString("grpclog-InfoDepth", fmt.Sprint(args...)).Fire()
	}
}

// WarningDepth logs to WARNING log at the specified depth. Arguments are handled in the manner of fmt.Print.
func (g *loggerT) WarningDepth(depth int, args ...interface{}) {
	if warningLevel < g.cfg.Level {
		return
	}
	_, file, line, ok := runtime.Caller(depth + 2)
	if ok {
		caller := fmt.Sprintf("%s:%d", file, line)
		g.output.Warn().RawString("grpclog-WarningDepth", fmt.Sprint(args...)).RawString("caller", caller).Fire()
	} else {
		g.output.Warn().RawString("grpclog-WarningDepth", fmt.Sprint(args...)).Fire()
	}
}

// ErrorDetph logs to ERROR log at the specified depth. Arguments are handled in the manner of fmt.Print.
func (g *loggerT) ErrorDepth(depth int, args ...interface{}) {
	if errorLevel < g.cfg.Level {
		return
	}
	_, file, line, ok := runtime.Caller(depth + 2)
	if ok {
		caller := fmt.Sprintf("%s:%d", file, line)
		g.output.Error().RawString("grpclog-ErrorDepth", fmt.Sprint(args...)).RawString("caller", caller).Fire()
	} else {
		g.output.Error().RawString("grpclog-ErrorDepth", fmt.Sprint(args...)).Fire()
	}
}

// FatalDepth logs to FATAL log at the specified depth. Arguments are handled in the manner of fmt.Print.
func (g *loggerT) FatalDepth(depth int, args ...interface{}) {
	if fatalLevel < g.cfg.Level {
		return
	}
	_, file, line, ok := runtime.Caller(depth + 2)
	if ok {
		caller := fmt.Sprintf("%s:%d", file, line)
		g.output.Fatal().RawString("grpclog-FatalDepth", fmt.Sprint(args...)).RawString("caller", caller).Fire()
	} else {
		g.output.Fatal().RawString("grpclog-FatalDepth", fmt.Sprint(args...)).Fire()
	}
	os.Exit(1)
}
