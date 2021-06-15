package kafkawrap

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/DataWorkbench/glog"
)

type ConsumerTrace struct {
	tracer opentracing.Tracer
	ctx    context.Context
}

func (ct *ConsumerTrace) OnConsume(msg *sarama.ConsumerMessage) {

	lp := glog.FromContext(ct.ctx)
	lp.Info().Msg("consumer interceptor OnConsume").Fire()

	//Extract messages headers
	mc := &MsgHeadersCarrier{}
	for _, h := range msg.Headers {
		mc.msgHeaders = append(mc.msgHeaders, *h)
	}
	producerSpan, err := ct.tracer.Extract(opentracing.TextMap, mc)
	if err != nil {
		lp.Error().Error("consumer interceptor extract error", err).Fire()
	}

	clientSpan := ct.tracer.StartSpan(
		"Kafka Consumer",
		opentracing.ChildOf(producerSpan),
		ext.SpanKindConsumer,
	)
	defer clientSpan.Finish()

}

// NewConsumerTrace processes span for intercepted messages
func NewConsumerTrace(ctx context.Context, tracer opentracing.Tracer) *ConsumerTrace {
	ct := ConsumerTrace{}
	ct.tracer = tracer
	ct.ctx = ctx
	return &ct
}
