package kafkawrap

import (
	"context"
	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
	"strings"
)

type KafkaConfig struct {
	Hosts    string `json:"hosts"         yaml:"hosts"         env:"HOSTS"             validate:"required"`
	GroupID  string `json:"group_id"      yaml:"group_id"      env:"GROUP_ID"          validate:"required"`
	Assignor string `json:"assignor"      yaml:"assignor"      env:"ASSIGNOR"          validate:"required"`
}

// NewProducerClient return a sarama kafka client,provide connection for kafka producer
func NewProducerClient(ctx context.Context, cfg *KafkaConfig, options ...Option) (client sarama.Client, tracer *ProducerTrace, err error) {

	opts := applyOptions(options...)
	lp := glog.FromContext(ctx)
	lp.Info().Msg("producer client connecting to kafka").String("hosts", cfg.Hosts).Fire()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	//opentracing jaeger producer interceptor
	tracer = NewProducerTrace(ctx, opts.tracer)
	config.Producer.Interceptors = []sarama.ProducerInterceptor{tracer}

	client, err = sarama.NewClient(strings.Split(cfg.Hosts, ","), config)

	return
}

// NewConsumerClient return a sarama kafka client,provide connection for kafka consumer
func NewConsumerClient(ctx context.Context, cfg *KafkaConfig, options ...Option) (client sarama.Client, err error) {

	opts := applyOptions(options...)
	lp := glog.FromContext(ctx)
	lp.Info().Msg("consumer client connecting to kafka").String("hosts", cfg.Hosts).Fire()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Interceptors = []sarama.ConsumerInterceptor{NewConsumerTrace(ctx, opts.tracer)}

	client, err = sarama.NewClient(strings.Split(cfg.Hosts, ","), config)

	return
}

// NewConsumerGroup return a sarama kafka consumer group
func NewConsumerGroup(ctx context.Context, cfg *KafkaConfig, options ...Option) (client sarama.ConsumerGroup, err error) {

	opts := applyOptions(options...)
	lp := glog.FromContext(ctx)
	lp.Info().Msg("consumer group client connecting to kafka").String("hosts", cfg.Hosts).Fire()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	switch cfg.Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		lp.Info().Msg("Unrecognized consumer group partition ").String("assignor", cfg.Assignor).Fire()
	}
	config.Consumer.Interceptors = []sarama.ConsumerInterceptor{NewConsumerTrace(ctx, opts.tracer)}

	client, err = sarama.NewConsumerGroup(strings.Split(cfg.Hosts, ","), cfg.GroupID, config)

	return
}
