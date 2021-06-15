package kafkawrap

import (
	"context"
	"log"
	"strings"

	"github.com/Shopify/sarama"

	"github.com/DataWorkbench/glog"
)

type ProducerConfig struct {
	Hosts string `json:"hosts"         yaml:"hosts"         env:"HOSTS"             validate:"required"`
}

type Producer struct {
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	tracer        *ProducerTrace
	lp            *glog.Logger
}

// NewProducerClient return a kafka producer
func NewProducerClient(ctx context.Context, cfg *ProducerConfig, options ...Option) (*Producer, error) {

	opts := applyOptions(options...)
	lp := glog.FromContext(ctx)
	lp.Info().Msg("producer client connecting to kafka").String("hosts", cfg.Hosts).Fire()
	tracer := NewProducerTrace(opts.tracer)

	s, err := newSyncProducer(tracer, cfg.Hosts)
	if err != nil {
		return nil, err
	}
	a, euro := newAsyncProducer(tracer, cfg.Hosts)
	if euro != nil {
		return nil, euro
	}

	p := &Producer{syncProducer: s, asyncProducer: a, tracer: tracer, lp: lp}

	return p, nil
}

func newSyncProducer(tracer *ProducerTrace, hosts string) (p sarama.SyncProducer, err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Interceptors = []sarama.ProducerInterceptor{tracer}
	p, err = sarama.NewSyncProducer(strings.Split(hosts, ","), config)
	return p, err
}

func newAsyncProducer(tracer *ProducerTrace, hosts string) (p sarama.AsyncProducer, err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = false
	config.Producer.Interceptors = []sarama.ProducerInterceptor{tracer}
	p, err = sarama.NewAsyncProducer(strings.Split(hosts, ","), config)
	return p, err
}

func (p *Producer) SyncProduce(ctx context.Context, topic string, msg []byte) (pid int32, offset int64, err error) {

	//provide context for opentracing.SpanFromContext
	p.tracer.ctx = ctx

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}
	pid, offset, err = p.syncProducer.SendMessage(message)

	return
}

func (p *Producer) AsyncProduce(ctx context.Context, topic string, msg []byte) {

	p.tracer.ctx = ctx

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}
	p.asyncProducer.Input() <- message
	// wait response
	select {
	case err := <-p.asyncProducer.Errors():
		log.Println("Produced message failure: ", err)
	default:
		log.Println("Produced message default")
	}
}

// Close wrapper for sarama producer Close()
func (p *Producer) Close() {
	if p == nil {
		return
	}
	p.lp.Info().Msg("waiting for sarama producer stop").Fire()
	if err := p.syncProducer.Close(); err != nil {
		p.lp.Error().Error("closing syncProducer error", err).Fire()
	}
	if err := p.asyncProducer.Close(); err != nil {
		p.lp.Error().Error("closing syncProducer error", err).Fire()
	}
	p.lp.Info().Msg("sarama producer stopped").Fire()
}
