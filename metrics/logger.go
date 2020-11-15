package metrics

import (
	"fmt"

	"github.com/DataWorkbench/glog"
)

// Logger for implements interface{} promhttp.Logger
type Logger struct {
	Output *glog.Logger
}

func (l *Logger) Println(v ...interface{}) {
	l.Output.Info().RawString("promhttp", fmt.Sprint(v...)).Fire()
}
