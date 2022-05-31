package logutil

import (
	"github.com/DataWorkbench/glog"
	"github.com/pkg/errors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

const (
	OutputConsole = "console"
	OutputFile    = "file"
	// console-file: console and file
	OutputConsoleFile = "console-file"
)


func New(cfg *Config) (*glog.Logger, error) {
	if err := checkConfig(cfg); err != nil {
		return nil, err
	}

	logger := glog.NewDefault().WithLevel(glog.Level(cfg.Level))
	switch cfg.Output {
	case OutputConsole:
	case OutputFile:
		logger = logger.WithExporter(glog.StandardExporter(fileWriter(*cfg)))
	case OutputConsoleFile:
		consoleExper := glog.StandardExporter(os.Stdout)
		fileExper := glog.StandardExporter(fileWriter(*cfg))
		logger = logger.WithExporter(glog.MultipleExporter(consoleExper, fileExper))
	}
	return logger, nil
}

type Config struct {
	// Level for set the log level. 1=>"debug", 2=>"info", 3=>"warn", 4=>"error", 5=>"fatal"
	Level int8 `json:"level"  yaml:"level" env:"LEVEL,default=1" validate:"gte=1,lte=5"`

	// Output specified the log output location. Optional value: "console" | "file" | "console-file"
	Output string `json:"output" yaml:"output" env:"OUTPUT,default=console" validate:"oneof=console file"`

	// File set the log file configuration. Is required if `Output` is "file"
	File *File `json:"file" yaml:"file" env:"FILE,default=" validate:"-"`
}

// File is configuration for log output to file.
type File struct {
	// Path is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Path string `json:"path" yaml:"path" env:"PATH" validate:"required"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"max_size" yaml:"max_size" env:"MAX_SIZE,default=128" validate:"gt=0"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"max_age" yaml:"max_age" env:"MAX_AGE,default=0" validate:"gte=0"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"max_backups" yaml:"max_backups" env:"MAX_BACKUPS,default=0" validate:"gte=0"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress" env:"COMPRESS" validate:"-"`
}

func checkConfig(cfg *Config) (err error) {
	if cfg == nil {
		return errors.New("logutil: the config cannot be nil")
	}
	if cfg.Level < 1 || cfg.Level > 5 {
		return errors.Errorf("logutil: the log level must greater than 0 and less than 6, you provided: %d", cfg.Level)
	}
	switch cfg.Output {
	case OutputConsole:
	case OutputFile, OutputConsoleFile:
		return checkFIleConfig(*cfg)
	default:
		return errors.New("logutil: log output must be oneof `console` or `file` or `console-file`")
	}
	return
}

func fileWriter(cfg Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.File.Path,
		MaxSize:    cfg.File.MaxSize,
		MaxAge:     cfg.File.MaxAge,
		MaxBackups: cfg.File.MaxBackups,
		LocalTime:  true,
		Compress:   cfg.File.Compress,
	}
}

func checkFIleConfig(cfg Config) error {
	if cfg.File == nil {
		return errors.New("logutil: the field file must be st if output file")
	}
	if cfg.File.Path == "" {
		return errors.New("logutil: the log file path must be set")
	}
	if cfg.File.MaxSize <= 0 {
		return errors.New("logutil: the max_size must be greater than 0")
	}
	if cfg.File.MaxAge < 0 {
		return errors.New("logutil: the max_age must be greater than or equal to 0")
	}
	if cfg.File.MaxBackups < 0 {
		return errors.New("logutil: the max_backups must be greater than or equal to 0")
	}
	return nil
}