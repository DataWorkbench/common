package kafka

import (
	"context"
	"regexp"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
)

// ConsumerWatcher used to consume regex format topics.
type ConsumerWatcher struct {
	group   *ConsumerGroup
	regexps map[string]*regexp.Regexp
}

// NewConsumerWatcher creates new ConsumerWatcher.
func NewConsumerWatcher(ctx context.Context, groupId string, cfg *ConsumerConfig, handler MessageHandler, options ...Option) (*ConsumerWatcher, error) {
	lp := glog.FromContext(ctx)

	group, err := NewConsumerGroup(ctx, groupId, cfg, handler, options...)
	if err != nil {
		return nil, err
	}

	c := &ConsumerWatcher{
		group:   group,
		regexps: make(map[string]*regexp.Regexp),
	}

	lp.Debug().Msg("ConsumerWatcher: successfully initialized consumer watcher").Fire()
	return c, nil
}

// Consume for start consumer in a loop.
// And monitor the topics changes that match regex.
//
// The loop will stop if the consumer closed or any unexpected errors happen.
//
// This function does not allow concurrent calls.
func (c *ConsumerWatcher) Consume(regexTopics []string) {
	c.initTopics(regexTopics)

	lg := c.group.lp
	lg.Debug().Msg("ConsumerWatcher: loop up and running").Strings("regexTopics", regexTopics).Fire()

	var err error
	var topics []string

LOOP:
	for {
		// CHeck if the consumer group was closed.
		select {
		case <-c.group.closed:
			break LOOP
		default:
		}

		// Fetch currently available topics.
		topics, err = c.fetchTopics()
		if err != nil {
			lg.Error().Error("ConsumerWatcher: fetch current topics error", err).Fire()
			if err == sarama.ErrClosedClient {
				break LOOP
			}
			time.Sleep(time.Second * 1)
			continue LOOP
		}

		roundCtx, roundCancel := context.WithCancel(c.group.ctx)
		go c.watchTopics(roundCtx, roundCancel, topics)

		// No qualified topics, waits for new topic joined..
		if len(topics) == 0 {
			lg.Debug().Msg("ConsumerWatcher: No qualified topics currently, wait for new topics to join").Fire()
			<-roundCtx.Done()
			continue LOOP
		}

		if err = c.group.consume(roundCtx, topics); err != nil {
			break LOOP
		}

		// CHeck if the consumer group was closed.
		select {
		case <-c.group.closed:
			break LOOP
		default:
		}

		// To prevent dead cycle.
		time.Sleep(time.Millisecond * 100)
	}

	lg.Debug().Msg("ConsumerWatcher: consumer was closed, stops").Strings("regexTopics", regexTopics).Fire()
}

// Close for close the consume group.
func (c *ConsumerWatcher) Close() (err error) {
	if c == nil {
		return
	}
	return c.group.Close()
}

func (c *ConsumerWatcher) initTopics(regexTopics []string) {
	for _, topic := range regexTopics {
		re := regexp.MustCompile(topic)
		c.regexps[topic] = re
	}
}

// fetchTopics to fetch the qualified topics.
func (c *ConsumerWatcher) fetchTopics() (topics []string, err error) {
	var availTopics []string
	// Get current available topics.
	availTopics, err = c.group.client.Topics()
	if err != nil {
		return
	}

	for _, topic := range availTopics {
		for _, re := range c.regexps {
			if re.MatchString(topic) {
				topics = append(topics, topic)
			}
		}
	}
	return
}

// watchTopics watching the subscription topics.
// This function exits and calls the `cancel` when the topics changed.
func (c *ConsumerWatcher) watchTopics(ctx context.Context, cancel context.CancelFunc, curTopics []string) {
	lg := c.group.lp
	ticker := time.NewTicker(c.group.client.Config().Metadata.RefreshFrequency)

	lg.Debug().Msg("ConsumerWatcher: watch for the topics changes").Fire()

	curMap := make(map[string]struct{}, len(curTopics))
	for _, topic := range curTopics {
		curMap[topic] = struct{}{}
	}

	var newTopics []string
	var err error
	var changed bool

LOOP:
	for {
		select {
		case <-ticker.C:
			newTopics, err = c.fetchTopics()
			if err != nil {
				if err == sarama.ErrClosedClient {
					lg.Error().Msg("ConsumerWatcher: client has been closed, exits...").Fire()
					break LOOP
				}
				lg.Error().Error("ConsumerWatcher: fetch new topics error", err).Fire()
				continue LOOP
			}

			if len(newTopics) != len(curTopics) {
				changed = true
				break LOOP
			}

			for _, topic := range newTopics {
				if _, ok := curMap[topic]; !ok {
					changed = true
					break LOOP
				}
			}
		case <-ctx.Done():
			lg.Debug().Msg("ConsumerWatcher: context was done, exits...").Fire()
			break LOOP
		}
	}

	if changed {
		lg.Debug().Msg("ConsumerWatcher: qualified topics changed").
			Strings("current", curTopics).
			Strings("new", newTopics).
			Fire()
	}

	ticker.Stop()
	cancel()
}
