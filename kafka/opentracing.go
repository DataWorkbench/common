package kafka

import (
	"context"

	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/DataWorkbench/common/gtrace"
)

var traceComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "sarama"}

type opentracingCarrier struct {
	headers []sarama.RecordHeader
}

// Set conforms to the opentracing.TextMapWriter interface.
func (c *opentracingCarrier) Set(key, val string) {
	rh := sarama.RecordHeader{Key: []byte(key), Value: []byte(val)}
	c.headers = append(c.headers, rh)
}

// ForeachKey conforms to the opentracing.TextMapReader interface.
func (c *opentracingCarrier) ForeachKey(handler func(key, val string) error) error {
	for i := range c.headers {
		h := c.headers[i]
		if err := handler(string(h.Key), string(h.Value)); err != nil {
			return err
		}
	}
	return nil
}

// producerTraceSpan start a span for producer.
func producerTraceSpan(ctx context.Context, tracer gtrace.Tracer, opName string) (span opentracing.Span, headers []sarama.RecordHeader) {
	if tid := gtrace.IdFromContext(ctx); tid != "" {
		headers = append(headers, sarama.RecordHeader{Key: []byte(gtrace.IdKey), Value: []byte(tid)})
	}

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span = tracer.StartSpan(
		opName,
		opentracing.ChildOf(parentCtx),
		ext.SpanKindProducer,
		traceComponentTag,
	)

	carrier := &opentracingCarrier{}
	err := tracer.Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		lg := glog.FromContext(ctx)
		lg.Error().Error("producerTraceSpan: tracer inject error", err).Fire()
	}

	headers = append(headers, carrier.headers...)
	return
}
