package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

// ClientConfig is the configuration for connects to kafka only as a client.
type ClientConfig struct {
	// The kafka hosts that split by `,`.  eg: "127.0.0.1:9092,127.0.0.1:9092"
	Hosts string `json:"hosts" yaml:"hosts" env:"HOSTS" validate:"required"`

	// RefreshFrequency is similar to `topic.metadata.refresh.interval.ms`
	// Defaults 10min.
	RefreshFrequency time.Duration `json:"refresh_frequency" yaml:"refresh_frequency" env:"REFRESH_FREQUENCY,default=10m" validate:"-"`
}

// convert the ConsumerConfig to sarama.Config
func (c *ClientConfig) convert() *sarama.Config {
	config := sarama.NewConfig()

	if c.RefreshFrequency <= 0 {
		config.Metadata.RefreshFrequency = time.Minute * 10 // defaults 10min.
	} else {
		config.Metadata.RefreshFrequency = c.RefreshFrequency
	}
	return config
}
