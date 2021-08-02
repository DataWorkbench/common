package kafka

import (
	"context"
	"strings"

	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"

	"github.com/DataWorkbench/common/trace"
)

// asyncProducer is wraps for sarama.AsyncProducer.
type asyncProducer struct {
	producer sarama.AsyncProducer
	lp       *glog.Logger
	tracer   opentracing.Tracer
}

// asyncMetadata used for pass-through data in AsyncProducer.
type asyncMetadata struct {
	span opentracing.Span
	tid  string // The trace id.
}

// NewAsyncProducer creates Producer with asyncProducer.
func NewAsyncProducer(ctx context.Context, cfg *ProducerConfig, options ...Option) (Producer, error) {
	opts := applyOptions(options...)
	lp := glog.FromContext(ctx)

	lp.Info().Msg("asyncProducer: initializing new async producer").String("hosts", cfg.Hosts).Fire()

	producer, err := sarama.NewAsyncProducer(strings.Split(cfg.Hosts, ","), cfg.convert())
	if err != nil {
		lp.Error().Error("asyncProducer: initializes async producer error", err).Fire()
		return nil, err
	}

	p := &asyncProducer{
		producer: producer,
		lp:       lp,
		tracer:   opts.tracer,
	}

	go p.checkSuccesses()
	go p.checkErrors()

	lp.Debug().Msg("asyncProducer: successfully initialized async producer").Fire()
	return p, nil
}

// Send sends message to kafka. The key allowed to be nil.
func (p *asyncProducer) Send(ctx context.Context, topic string, key Encoder, value Encoder) (err error) {
	span, headers := producerTraceSpan(ctx, p.tracer, "AsyncProduceMessage")

	message := &sarama.ProducerMessage{
		Topic:    topic,
		Key:      key,
		Value:    value,
		Headers:  headers,
		Metadata: &asyncMetadata{span: span, tid: trace.IdFromContext(ctx)},
	}

	p.producer.Input() <- message
	return
}

// Close close the AsyncProducer.
func (p *asyncProducer) Close() (err error) {
	if p == nil {
		return
	}

	p.lp.Debug().Msg("asyncProducer: wait for the producer to close").Fire()
	if err = p.producer.Close(); err != nil {
		p.lp.Error().Error("asyncProducer: producer close error", err).Fire()
		return
	}
	p.lp.Debug().Msg("asyncProducer: producer successful closed").Fire()
	return
}

func (p *asyncProducer) checkSuccesses() {
	for msg := range p.producer.Successes() {
		meta := msg.Metadata.(*asyncMetadata)

		p.lp.Debug().Msg("asyncProducer: send message success").
			String("topic", msg.Topic).
			Int32("partition", msg.Partition).
			Int64("offset", msg.Offset).
			String(trace.IdKey, meta.tid).
			Fire()

		span := meta.span

		span.SetTag("topic", msg.Topic)
		span.SetTag("partition", msg.Partition)
		span.SetTag("offset", msg.Offset)

		span.Finish()
	}

	p.lp.Debug().Msg("asyncProducer.checkSuccesses: channel has been closed, exits").Fire()
}

func (p *asyncProducer) checkErrors() {
	for pe := range p.producer.Errors() {
		msg := pe.Msg
		meta := pe.Msg.Metadata.(*asyncMetadata)

		p.lp.Error().Msg("asyncProducer: send message failed").
			String("topic", msg.Topic).
			Error("error", pe.Err).
			String(trace.IdKey, meta.tid).
			Fire()

		span := meta.span

		span.SetTag("topic", msg.Topic)
		span.SetTag("partition", msg.Partition)
		span.SetTag("offset", msg.Offset)

		ext.Error.Set(span, true)
		span.LogFields(tracerLog.Error(pe.Err))

		span.Finish()
	}

	p.lp.Debug().Msg("asyncProducer.checkErrors: channel has been closed, exits").Fire()
}
