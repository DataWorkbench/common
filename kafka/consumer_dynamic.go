package kafka

import (
	"context"
	"sync"
	"time"
)

// ConsumerDynamic used to consume the topics with dynamic.
type ConsumerDynamic struct {
	group *ConsumerGroup

	carrier chan []string
	topics  []string

	mux *sync.Mutex // protects access to the topics.
}

// NewConsumerDynamic creates new ConsumerDynamic.
func NewConsumerDynamic(ctx context.Context, groupId string, cfg *ConsumerConfig, handler MessageHandler, options ...Option) (*ConsumerDynamic, error) {
	group, err := NewConsumerGroup(ctx, groupId, cfg, handler, options...)
	if err != nil {
		return nil, err
	}

	c := &ConsumerDynamic{
		group:   group,
		carrier: make(chan []string),
		topics:  nil,
		mux:     new(sync.Mutex),
	}
	group.lp.Debug().Msg("ConsumerDynamic: successfully initialized consumer dynamic").Fire()
	return c, nil
}

func (c *ConsumerDynamic) setTopics(topics []string) {
	c.mux.Lock()
	c.topics = topics
	c.mux.Unlock()
}

func (c *ConsumerDynamic) getTopics() (topics []string) {
	c.mux.Lock()
	topics = c.topics
	c.mux.Unlock()
	return
}

func (c *ConsumerDynamic) topicHandler(ctx context.Context, topics []string, _ []string, _ []string) error {
	select {
	case c.carrier <- topics:
	case <-c.group.closed:
	case <-ctx.Done():
	}
	return nil
}

func (c *ConsumerDynamic) watchTopics(ctx context.Context, cancel context.CancelFunc) {
	lg := c.group.lp
	lg.Debug().Msg("ConsumerDynamic: watch for the topics changes").Fire()

	select {
	case topics := <-c.carrier:
		lg.Debug().Msg("ConsumerDynamic: qualified topics changed").
			Strings("current", c.getTopics()).
			Strings("new", topics).
			Fire()
		c.setTopics(topics)
	case <-c.group.closed:
	case <-ctx.Done():
	}
	cancel()
}

// Consume for start consumer in a loop.
// And monitor the topics changes that match regex.
// regexTopics eg: ["^a-.*$", "^b-$"]
//
// The loop will stop if the consumer closed or any unexpected errors happen.
//
// This function does not allow concurrent calls.
func (c *ConsumerDynamic) Consume(regexTopics []string) (err error) {
	c.group.wg.Add(1)
	defer c.group.wg.Done()

	lg := c.group.lp

	lg.Debug().Msg("ConsumerDynamic: loop up and running").Strings("regexTopics", regexTopics).Fire()

	var watcher *TopicWatcher
	watcher, err = newTopicWatcherFromClient(c.group.ctx, c.group.client, regexTopics, c.topicHandler)
	if err != nil {
		return
	}

	go func() {
		_ = watcher.Watch()
	}()

	defer func() {
		lg.Debug().Msg("ConsumerDynamic: consumer was closed, stops").Fire()
		_ = watcher.Close()
	}()

	lg.Debug().Msg("ConsumerDynamic: wait for the topics to join").Fire()

	// start watch topics
	select {
	case topics := <-c.carrier:
		c.topics = topics
	case <-c.group.ctx.Done():
		return
	case <-c.group.closed:
		return
	}

LOOP:
	for {
		// CHeck if the consumer group was closed.
		select {
		case <-c.group.ctx.Done():
			break LOOP
		case <-c.group.closed:
			return
		default:
		}

		roundCtx, roundCancel := context.WithCancel(c.group.ctx)
		go c.watchTopics(roundCtx, roundCancel)

		// No qualified topics, waits for new topic joined..
		if len(c.getTopics()) == 0 {
			lg.Debug().Msg("ConsumerDynamic: No qualified topics currently").Fire()
			<-roundCtx.Done()
			continue LOOP
		}

		if err = c.group.consume(roundCtx, c.getTopics()); err != nil {
			break LOOP
		}

		// CHeck if the consumer group was closed.
		select {
		case <-c.group.ctx.Done():
			break LOOP
		case <-c.group.closed:
			return
		default:
		}

		lg.Debug().Msg("ConsumerDynamic: re-balance happens, topics or partitions or consumers changed").Fire()

		// To prevent dead cycle.
		time.Sleep(time.Millisecond * 100)
	}
	return
}

// Close for close the consume group.
func (c *ConsumerDynamic) Close() (err error) {
	if c == nil {
		return
	}
	return c.group.Close()
}
