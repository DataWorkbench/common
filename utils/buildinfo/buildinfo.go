package buildinfo

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultValue = "unknown"
)

// Version information, compile-time settings by go build -ldflags.
// Reference the build.sh
var (
	GoVersion   = defaultValue
	CompileBy   = defaultValue
	CompileTime = defaultValue
	GitBranch   = defaultValue
	GitCommit   = defaultValue
	OsArch      = defaultValue
)

var (
	// MapValue return build info in map type value
	MapValue map[string]string

	// JSONString is a strings in JSON format.
	// e.g:
	// {"go_version":"go1.15.3","compile_by":"xxxxxx@yunify.com","compile_time":"2020-11-15:17:02:52","git_branch":"dev","git_commit":"7fcbe3f","os_arch":"Darwin/x86_64"}
	JSONString string

	// SingleString is a strings in single line format.
	// e.g: go_version=go1.15.3 compile_by=xxxxxx@yunify.com compile_time=2020-11-15:17:02:52 git_branch=dev git_commit=7fcbe3f os_arch=Darwin/x86_64
	SingleString string

	// MultiString is a strings in multiple line format.
	// e.g:
	/*
	   compile_time    2020-11-15:17:02:52
	   git_branch      dev
	   git_commit      7fcbe3f
	   os_arch         Darwin/x86_64
	   go_version      go1.15.3
	   compile_by      xxxxxx@yunify.com
	*/
	MultiString string
)

func init() {
	MapValue = buildMapValue()
	JSONString = buildJSONString()
	SingleString = buildSingleString()
	MultiString = buildMultiString()

	// register prometheus metrics
	prometheus.MustRegister(createCollector())
}

func buildMapValue() map[string]string {
	return map[string]string{
		"go_version":   GoVersion,
		"compile_by":   CompileBy,
		"compile_time": CompileTime,
		"git_branch":   GitBranch,
		"git_commit":   GitCommit,
		"os_arch":      OsArch,
	}
}

// createCollector will create a gauge metric named `program_build_info`,
// which include multiple labels to expose server's build info.
func createCollector() prometheus.Collector {
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   "program",
			Name:        "build_info",
			Help:        "A metric with a constant '1' value labeled by program build info.",
			ConstLabels: MapValue,
		},
		func() float64 { return 1 },
	)
}

func buildJSONString() string {
	var b strings.Builder

	m := MapValue
	n := len(m) - 1
	i := 0

	b.WriteByte('{')
	for k, v := range m {
		b.WriteByte('"')
		b.WriteString(k)
		b.WriteByte('"')

		b.WriteByte(':')

		b.WriteByte('"')
		b.WriteString(v)
		b.WriteByte('"')

		if i != n {
			b.WriteByte(',')
			i++
		}
	}
	b.WriteByte('}')
	return b.String()
}

func buildSingleString() string {
	var b strings.Builder

	m := MapValue
	n := len(m) - 1
	i := 0

	for k, v := range m {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)

		if i != n {
			b.WriteByte(' ')
			i++
		}
	}
	return b.String()
}

func buildMultiString() string {
	var b strings.Builder

	m := MapValue
	n := len(m) - 1
	i := 0

	for k, v := range m {
		b.WriteString(fmt.Sprintf("%-15s %s", k, v))

		if i != n {
			b.WriteByte('\n')
			i++
		}
	}
	return b.String()
}
