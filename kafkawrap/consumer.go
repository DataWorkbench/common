package kafkawrap

import (
	"context"
	"strings"

	"github.com/Shopify/sarama"

	"github.com/DataWorkbench/glog"
)

type ConsumerConfig struct {
	Hosts    string `json:"hosts"         yaml:"hosts"         env:"HOSTS"             validate:"required"`
	GroupID  string `json:"group_id"      yaml:"group_id"      env:"GROUP_ID"          validate:"required"`
	Assignor string `json:"assignor"      yaml:"assignor"      env:"ASSIGNOR"          validate:"required"`
}

type ConsumerGroup struct {
	Group  sarama.ConsumerGroup
	cb     Callback
	ready  chan bool
	lp     *glog.Logger
	ctx    context.Context
	cancel func()
}

// NewConsumerGroup return a sarama kafka consumer group
func NewConsumerGroup(ctx context.Context, cfg *ConsumerConfig, options ...Option) (*ConsumerGroup, error) {

	opts := applyOptions(options...)
	lp := glog.FromContext(ctx)
	lp.Info().Msg("consumer group client connecting to kafka").String("hosts", cfg.Hosts).Fire()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
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

	client, err := sarama.NewConsumerGroup(strings.Split(cfg.Hosts, ","), cfg.GroupID, config)

	consumer := &ConsumerGroup{Group: client, lp: lp}
	return consumer, err
}

//eg: consumer.ConsumeWithGroup(ctx, "topicE",func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
//	   for message := range claim.Messages() {
//  	 	log.Println("consumer message ",string(message.Value),"Offset",message.Offset,"Partition",message.Partition)
//	    	session.MarkMessage(message, "")
//	   }
//    })

func (c *ConsumerGroup) ConsumeWithGroup(ctx context.Context, topic string, cb Callback) {
	if c == nil {
		return
	}
	c.cb = cb
	c.ready = make(chan bool)
	c.ctx, c.cancel = context.WithCancel(ctx)
	go func() {
		for {
			// sarama.ConsumerGroup.Consume() should be called inside an infinite loop,
			// when rebalance happens, the consumer session will need to be recreated to get the new claims
			if err := c.Group.Consume(c.ctx, []string{topic}, c); err != nil {
				c.lp.Error().Error("consumer group error", err).Fire()
			}
			c.lp.Info().Msg("partition rebalance happens, a claim session exit, create a new consumer session").Fire()

			// check if context was cancelled, signaling that the consumer should stop
			if c.ctx.Err() != nil {
				c.lp.Info().Msg("context was cancelled,stop consumer").Fire()
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	c.lp.Info().Msg("Sarama consumer up and running!...").Fire()
}

type Callback func(sess sarama.ConsumerGroupSession, cc sarama.ConsumerGroupClaim)

func (c *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	// mark every consumer session as ready and avoid deadlock
	close(c.ready)
	return nil
}

func (c *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.cb(session, claim)
	return nil
}

// Close wrapper for sarama.ConsumerGroup.Close(),Call when exit the app
func (c *ConsumerGroup) Close() {
	if c == nil {
		return
	}
	// end the sarama.ConsumerGroup.Consume() loop before Close()
	c.cancel()

	c.lp.Info().Msg("waiting for sarama consumer group stop").Fire()
	// stops the ConsumerGroup and detaches any running sessions
	if err := c.Group.Close(); err != nil {
		c.lp.Error().Error("closing client error", err).Fire()
	}
	c.lp.Info().Msg("sarama consumer group stopped").Fire()
}
