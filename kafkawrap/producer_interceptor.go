package kafkawrap

import (
	"context"
	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type ProducerTrace struct {
	tracer opentracing.Tracer
	Ctx    context.Context
}

func (pt *ProducerTrace) OnSend(msg *sarama.ProducerMessage) {

	lp := glog.FromContext(pt.Ctx)
	lp.Info().Msg("kafka producer interceptor OnSend").Fire()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(pt.Ctx); parent != nil {
		parentCtx = parent.Context()
	}

	//tags := opentracing.Tags{
	//	"kafka.message.topic":  topicsStr,
	//	"kafka.message.length": length,
	//	"span.kind":            "client",
	//}
	clientSpan := pt.tracer.StartSpan(
		"Kafka Producer",
		opentracing.ChildOf(parentCtx),
		ext.SpanKindProducer, //tags,
	)

	mc := &MsgHeadersCarrier{}
	err := pt.tracer.Inject(clientSpan.Context(), opentracing.TextMap, mc)
	if err != nil {
		lp.Error().Error("producer interceptor inject failed", err).Fire()
	}
	msg.Headers = append(msg.Headers, mc.msgHeaders...)

	defer clientSpan.Finish()

}

// NewProducerTrace processes add some headers with the span data.
func NewProducerTrace(ctx context.Context, tracer opentracing.Tracer) *ProducerTrace {
	pt := ProducerTrace{}
	pt.tracer = tracer
	pt.Ctx = ctx
	return &pt
}

type MsgHeadersCarrier struct {
	msgHeaders []sarama.RecordHeader
}

// Set conforms to the TextMapWriter interface.
func (c *MsgHeadersCarrier) Set(key, val string) {
	rh := sarama.RecordHeader{Key: []byte(key), Value: []byte(val)}
	c.msgHeaders = append(c.msgHeaders, rh)
}

// ForeachKey conforms to the TextMapReader interface.
func (c *MsgHeadersCarrier) ForeachKey(handler func(key, val string) error) error {
	for _, h := range c.msgHeaders {
		if err := handler(string(h.Key), string(h.Value)); err != nil {
			return err
		}
	}
	return nil
}
