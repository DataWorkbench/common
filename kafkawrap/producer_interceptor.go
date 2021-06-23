package kafkawrap

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"

	"github.com/DataWorkbench/glog"
)

type producerTrace struct {
	tracer opentracing.Tracer
	ctx    context.Context
}

func (pt *producerTrace) OnSend(msg *sarama.ProducerMessage) {

	lp := glog.FromContext(msg.Metadata.(context.Context))
	lp.Info().Msg("kafka producer interceptor OnSend").Fire()

	//lp.WithFields().AddString(ctxReqIdKey, reqId)  暂时没找到lp get方法，先从grpc metadata获取rid
	md, ok := metadata.FromIncomingContext(msg.Metadata.(context.Context))
	var ids []string
	if ok {
		ids = md.Get("rid")
	}

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(msg.Metadata.(context.Context)); parent != nil {
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

	msg.Headers = append(msg.Headers, sarama.RecordHeader{Key: []byte("rid"), Value: []byte(ids[0])})

	clientSpan.Finish()

}

// NewProducerTrace processes add some headers with the span data.
func NewProducerTrace(tracer opentracing.Tracer) *producerTrace {
	pt := producerTrace{}
	pt.tracer = tracer
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
