package gtrace

import (
	"io"

	"github.com/DataWorkbench/glog"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type Tracer = opentracing.Tracer
type SpanContext = jaeger.SpanContext

// Config for create a new tracer.
type Config struct {
	ServiceName string `json:"service_name" yaml:"service_name" env:"service_name" validate:"required"`
	LocalAgent  string `json:"local_agent"  yaml:"local_agent"  env:"local_agent"  validate:"required"`
}

// New create a new opentracing.Tracer by jaeger.
func New(cfg *Config) (tracer Tracer, closer io.Closer, err error) {
	// Config the jaeger
	jCfg := genJaegerConfig(cfg)
	tracer, closer, err = jCfg.NewTracer(config.Logger(jaeger.NullLogger))
	return
}

func NewWithGLog(cfg *Config, lg *glog.Logger) (tracer Tracer, closer io.Closer, err error) {
	// Config the jaeger
	jCfg := genJaegerConfig(cfg)
	tracer, closer, err = jCfg.NewTracer(config.Logger(&logger{Output: lg}))
	return
}

func genJaegerConfig(cfg *Config) config.Configuration {
	return config.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		// Set agent collect
		Reporter: &config.ReporterConfig{
			// 127.0.0.1:6831
			LocalAgentHostPort: cfg.LocalAgent,
		},

		// default is uber-trace-id
		Headers: &jaeger.HeadersConfig{
			JaegerDebugHeader:        "x-trace-debug-id",
			JaegerBaggageHeader:      "x-trace-baggage",
			TraceContextHeaderName:   "x-trace-id",
			TraceBaggageHeaderPrefix: "x-trace-ctx",
		},
	}
}
