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

type msgHeadersCarrier struct {
	msgHeaders []sarama.RecordHeader
}

// Set conforms to the TextMapWriter interface.
func (c *msgHeadersCarrier) Set(key, val string) {
	rh := sarama.RecordHeader{Key: []byte(key), Value: []byte(val)}
	c.msgHeaders = append(c.msgHeaders, rh)
}

// ForeachKey conforms to the TextMapReader interface.
func (c *msgHeadersCarrier) ForeachKey(handler func(key, val string) error) error {
	for i := range c.msgHeaders {
		h := c.msgHeaders[i]
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

	mc := &msgHeadersCarrier{}
	err := tracer.Inject(span.Context(), opentracing.TextMap, mc)
	if err != nil {
		lg := glog.FromContext(ctx)
		lg.Error().Error("producerTraceSpan: tracer inject error", err).Fire()
	}

	headers = append(headers, mc.msgHeaders...)
	return
}
