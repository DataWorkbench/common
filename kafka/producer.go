package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

var (
	_ Producer = (*syncProducer)(nil)
	_ Producer = (*asyncProducer)(nil)
)

// type helpful for caller reference.
type (
	Encoder       = sarama.Encoder
	ByteEncoder   = sarama.ByteEncoder
	StringEncoder = sarama.StringEncoder
)

type Producer interface {
	Send(ctx context.Context, topic string, key Encoder, value Encoder) error
	Close() error
}

// ProducerConfig is the configuration for connects to kafka as a producer.
type ProducerConfig struct {
	// The kafka hosts that split by `,`. eg: "127.0.0.1:9092,127.0.0.1:9092"
	Hosts string `json:"hosts" yaml:"hosts" env:"HOSTS" validate:"required"`

	// RequiredAcks is similar as request.required.acks.
	// Optional values: 0 => NoResponse, 1 => WaitForLocal, -1 => WaitForAll.
	// Defaults -1.
	RequiredAcks int `json:"required_acks" yaml:"required_acks" env:"REQUIRED_ACKS,default=-1" validate:"oneof=0 1 -1"`

	// PartitionerClass is similar as partitioner.class
	// Optional values: "hash", "random", "roundRobin", "manual", "referenceHash"
	// Defaults "hash".
	PartitionerClass string `json:"partitioner_class" yaml:"partitioner_class" env:"PARTITIONER_CLASS,default=hash" validate:"oneof=hash random roundRobin manual referenceHash"`
}

// convert the ProducerConfig to sarama.Config
func (c *ProducerConfig) convert() *sarama.Config {
	config := sarama.NewConfig()

	config.MetricRegistry = metricRegistry
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	switch c.RequiredAcks {
	case 0:
		config.Producer.RequiredAcks = sarama.NoResponse
	case 1:
		config.Producer.RequiredAcks = sarama.WaitForLocal
	case -1:
		config.Producer.RequiredAcks = sarama.WaitForAll
	default:
		config.Producer.RequiredAcks = sarama.WaitForAll
	}

	switch c.PartitionerClass {
	case "hash":
		config.Producer.Partitioner = sarama.NewHashPartitioner
	case "random":
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	case "roundRobin":
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case "manual":
		config.Producer.Partitioner = sarama.NewManualPartitioner
	case "referenceHash":
		config.Producer.Partitioner = sarama.NewReferenceHashPartitioner
	default:
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}

	return config
}
