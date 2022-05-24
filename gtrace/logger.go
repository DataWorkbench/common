package gtrace

import (
	"fmt"

	"github.com/DataWorkbench/glog"
)

type logger struct {
	Output *glog.Logger // the default logger instances
}

// Error logs a message at error priority
func (l *logger) Error(msg string) {
	l.Output.Error().RawString("gtrace", msg).Fire()
}

// Infof logs a message at info priority
func (l *logger) Infof(msg string, args ...interface{}) {
	l.Output.Info().RawString("gtrace", fmt.Sprintf(msg, args...)).Fire()
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	l.Output.Debug().RawString("gtrace", fmt.Sprintf(msg, args...)).Fire()
}
