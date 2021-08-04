package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

// ConsumerConfig is the configuration for connects to kafka as a consumer.
type ConsumerConfig struct {
	// The kafka hosts that split by `,`.  eg: "127.0.0.1:9092,127.0.0.1:9092"
	Hosts string `json:"hosts" yaml:"hosts" env:"HOSTS" validate:"required"`

	// RefreshFrequency is similar to `topic.metadata.refresh.interval.ms`
	// Defaults 10min.
	RefreshFrequency time.Duration `json:"refresh_frequency" yaml:"refresh_frequency" env:"REFRESH_FREQUENCY,default=10m" validate:"-"`

	// Consumer.Offsets.Initial. -1 => OffsetNewest, -2 => OffsetOldest
	// Defaults OffsetOldest.
	OffsetsInitial int64 `json:"offsets_initial" yaml:"offsets_initial" env:"OFFSETS_INITIAL,default=-2" validate:"oneof=-2 -1"`

	// Consumer.Group.Rebalance.Strategy.
	// Optional values: "sticky", "range", "roundRobin".
	// Defaults "roundRobin".
	BalanceStrategy string `json:"balance_strategy" yaml:"balance_strategy" env:"BALANCE_STRATEGY,default=sticky" validate:"oneof=sticky range roundRobin"`
}

// convert the ConsumerConfig to sarama.Config
func (c *ConsumerConfig) convert() *sarama.Config {
	config := sarama.NewConfig()

	config.MetricRegistry = metricRegistry
	config.Consumer.Return.Errors = true

	// Sets follows parameters by `ConsumerConfig` if necessary.
	config.Consumer.Group.Session.Timeout = time.Second * 10
	config.Consumer.Group.Heartbeat.Interval = time.Second * 3 // By defaults.

	if c.RefreshFrequency <= 0 {
		config.Metadata.RefreshFrequency = time.Minute * 10 // defaults 10min.
	} else {
		config.Metadata.RefreshFrequency = c.RefreshFrequency
	}

	switch c.OffsetsInitial {
	case -1:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	case -2:
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	default:
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	switch c.BalanceStrategy {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	case "roundRobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	default:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	}

	return config
}
