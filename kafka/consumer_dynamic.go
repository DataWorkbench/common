package kafka

import (
	"context"
	"regexp"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
)

// ConsumerDynamic used to consume the topics with dynamic.
type ConsumerDynamic struct {
	group   *ConsumerGroup
	regexps map[string]*regexp.Regexp
}

// NewConsumerDynamic creates new ConsumerDynamic.
func NewConsumerDynamic(ctx context.Context, groupId string, cfg *ConsumerConfig, handler MessageHandler, options ...Option) (*ConsumerDynamic, error) {
	lp := glog.FromContext(ctx)

	group, err := NewConsumerGroup(ctx, groupId, cfg, handler, options...)
	if err != nil {
		return nil, err
	}

	c := &ConsumerDynamic{
		group:   group,
		regexps: make(map[string]*regexp.Regexp),
	}

	lp.Debug().Msg("ConsumerDynamic: successfully initialized consumer watcher").Fire()
	return c, nil
}

// Consume for start consumer in a loop.
// And monitor the topics changes that match regex.
// regexTopics eg: ["^a-.*$", "^b-$"]
//
// The loop will stop if the consumer closed or any unexpected errors happen.
//
// This function does not allow concurrent calls.
func (c *ConsumerDynamic) Consume(regexTopics []string) {
	c.initTopics(regexTopics)

	lg := c.group.lp
	lg.Debug().Msg("ConsumerDynamic: loop up and running").Strings("regexTopics", regexTopics).Fire()

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
			lg.Error().Error("ConsumerDynamic: fetch current topics error", err).Fire()
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
			lg.Debug().Msg("ConsumerDynamic: No qualified topics currently, wait for new topics to join").Fire()
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

	lg.Debug().Msg("ConsumerDynamic: consumer was closed, stops").Strings("regexTopics", regexTopics).Fire()
}

// Close for close the consume group.
func (c *ConsumerDynamic) Close() (err error) {
	if c == nil {
		return
	}
	return c.group.Close()
}

func (c *ConsumerDynamic) initTopics(regexTopics []string) {
	for _, topic := range regexTopics {
		re := regexp.MustCompile(topic)
		c.regexps[topic] = re
	}
}

// fetchTopics to fetch the qualified topics.
func (c *ConsumerDynamic) fetchTopics() (topics []string, err error) {
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
func (c *ConsumerDynamic) watchTopics(ctx context.Context, cancel context.CancelFunc, curTopics []string) {
	lg := c.group.lp
	ticker := time.NewTicker(c.group.client.Config().Metadata.RefreshFrequency)

	lg.Debug().Msg("ConsumerDynamic: watch for the topics changes").Fire()

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
					lg.Error().Msg("ConsumerDynamic: client has been closed, exits...").Fire()
					break LOOP
				}
				lg.Error().Error("ConsumerDynamic: fetch new topics error", err).Fire()
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
			lg.Debug().Msg("ConsumerDynamic: context was done, exits...").Fire()
			break LOOP
		}
	}

	if changed {
		lg.Debug().Msg("ConsumerDynamic: qualified topics changed").
			Strings("current", curTopics).
			Strings("new", newTopics).
			Fire()
	}

	ticker.Stop()
	cancel()
}
