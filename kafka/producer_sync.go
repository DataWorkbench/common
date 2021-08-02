package kafka

import (
	"context"
	"strings"

	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
)

// syncProducer is wraps for sarama.SyncProducer.
type syncProducer struct {
	producer sarama.SyncProducer
	lp       *glog.Logger
	tracer   opentracing.Tracer
}

// NewSyncProducer creates Producer with syncProducer.
func NewSyncProducer(ctx context.Context, cfg *ProducerConfig, options ...Option) (Producer, error) {
	opts := applyOptions(options...)
	lp := glog.FromContext(ctx)

	lp.Info().Msg("syncProducer: initializing new sync producer").String("hosts", cfg.Hosts).Fire()

	producer, err := sarama.NewSyncProducer(strings.Split(cfg.Hosts, ","), cfg.convert())
	if err != nil {
		lp.Error().Error("syncProducer: initializes sync producer error", err).Fire()
		return nil, err
	}

	p := &syncProducer{
		producer: producer,
		lp:       lp,
		tracer:   opts.tracer,
	}

	lp.Debug().Msg("syncProducer: successfully initialized sync producer").Fire()
	return p, nil
}

// Send sends message to kafka. The key allowed to be nil.
func (p *syncProducer) Send(ctx context.Context, topic string, key Encoder, value Encoder) (err error) {
	var partition int32
	var offset int64

	lg := glog.FromContext(ctx)
	span, headers := producerTraceSpan(ctx, p.tracer, "SyncProduceMessage")

	message := &sarama.ProducerMessage{
		Topic:    topic,
		Key:      key,
		Value:    value,
		Headers:  headers,
		Metadata: nil,
	}

	partition, offset, err = p.producer.SendMessage(message)
	if err != nil {
		lg.Error().Msg("syncProducer: send message failed").
			String("topic", topic).
			Error("error", err).
			Fire()
	} else {
		lg.Debug().Msg("syncProducer: send message success").
			String("topic", topic).
			Int32("partition", partition).
			Int64("offset", offset).
			Fire()
	}

	span.SetTag("topic", topic)
	span.SetTag("partition", partition)
	span.SetTag("offset", offset)

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(tracerLog.Error(err))
	}

	// Finish the opentracing span.
	span.Finish()
	return
}

// Close close the SyncProducer.
func (p *syncProducer) Close() (err error) {
	if p == nil {
		return
	}

	p.lp.Debug().Msg("syncProducer: wait for the producer to close").Fire()
	if err = p.producer.Close(); err != nil {
		p.lp.Error().Error("syncProducer: producer close error", err).Fire()
		return
	}
	p.lp.Debug().Msg("syncProducer: producer successful closed").Fire()
	return
}
