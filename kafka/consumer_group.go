package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
)

// ConsumerGroup is wraps for sarama.ConsumerGroup.
type ConsumerGroup struct {
	ctx    context.Context
	lp     *glog.Logger
	client sarama.Client
	group  sarama.ConsumerGroup

	// Initialize by inside.
	handler sarama.ConsumerGroupHandler

	closed chan struct{}
}

// NewConsumerGroup creates a new ConsumerGroup.
func NewConsumerGroup(ctx context.Context, groupId string, cfg *ConsumerConfig, handler MessageHandler, options ...Option) (*ConsumerGroup, error) {
	if groupId == "" {
		panic("ConsumerGroup: groupId can not be empty")
	}

	lp := glog.FromContext(ctx)

	lp.Info().Msg("ConsumerGroup: initializing new kafka client").String("hosts", cfg.Hosts).Fire()

	client, err := sarama.NewClient(strings.Split(cfg.Hosts, ","), cfg.convert())
	if err != nil {
		lp.Error().Error("ConsumerGroup: initializes kafka client error", err).Fire()
		return nil, err
	}

	group, err := sarama.NewConsumerGroupFromClient(groupId, client)
	if err != nil {
		lp.Error().Error("ConsumerGroup: initializes consumer cfg error", err).Fire()
		return nil, err
	}

	c := &ConsumerGroup{
		ctx:     ctx,
		lp:      lp,
		client:  client,
		group:   group,
		handler: newConsumerHandler(ctx, handler, options...),
		closed:  make(chan struct{}),
	}

	lp.Debug().Msg("ConsumerGroup: successfully initialized consumer group").Fire()
	return c, nil
}

func (c *ConsumerGroup) consume(ctx context.Context, topics []string) (err error) {
	lg := c.lp

	lg.Debug().Msg("ConsumerGroup: group consume started").Strings("topics", topics).Fire()

	// when re-balance happens, the consume session will need to be recreated to get the new claims.
	if err = c.group.Consume(ctx, topics, c.handler); err != nil {
		lg.Error().Error("ConsumerGroup: group consume error", err).Strings("topics", topics).Fire()
		return
	}

	// process the consumer errors.
LOOP:
	for {
		select {
		case err, ok := <-c.group.Errors():
			if !ok {
				break LOOP
			}
			if pe, ok := err.(*sarama.ConsumerError); ok {
				if pe.Err == context.Canceled {
					continue LOOP
				}
			}
			lg.Error().Error("ConsumerGroup: consumer returns error", err).Strings("topics", topics).Fire()
		default:
			break LOOP
		}
	}
	return
}

// Consume for start consumer in a loop and automatic processing re-balance.
//
// The loop will stop if the consumer closed or any unexpected errors happen.
//
// This function does not allow concurrent calls.
func (c *ConsumerGroup) Consume(topics []string) {
	lg := c.lp

	lg.Debug().Msg("ConsumerGroup: loop up and running").Strings("topics", topics).Fire()

LOOP:
	for {
		if err := c.consume(c.ctx, topics); err != nil {
			break LOOP
		}

		// Check if the consumer group was closed.
		select {
		case <-c.closed:
			break LOOP
		default:
		}

		lg.Debug().Msg("ConsumerGroup: re-balance happens, topics or partitions or consumers changed").Strings("topics", topics).Fire()

		// To prevent dead cycle.
		time.Sleep(time.Millisecond * 100)
	}

	lg.Debug().Msg("ConsumerGroup: consumer was closed, stops").Strings("topics", topics).Fire()
}

// Close wrapper for sarama.ConsumerGroup.Close(), Calls before exit the app.
func (c *ConsumerGroup) Close() (err error) {
	if c == nil {
		return
	}
	close(c.closed)

	c.lp.Debug().Msg("ConsumerGroup: wait for the consumer to close").Fire()
	// stops the ConsumerGroup and detaches any running sessions
	if err = c.group.Close(); err != nil {
		c.lp.Error().Error("ConsumerGroup: close consumer error", err).Fire()
		return
	}
	c.lp.Debug().Msg("ConsumerGroup: consumer successful closed").Fire()
	return
}
