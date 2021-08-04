package kafka

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
)

// TopicHandler called when the list of topics changed.
//
// `TopicWatcher` will block until all handler exit before closing. Thus, you must be monitored `<- ctx.Done` to
// receive the exit signal if you implementation is resident.
//
// The function parameter is:
//  - ctx: The `ctx` is created with `context.WithCancel`.
//  - wg:
//  - topics: The list of topics that currently qualified.
//  - increases: The list of topics added compared to the previous cycle.
//  - decreases: The list of topics reduced compared to the previous cycle
type TopicHandler func(ctx context.Context, wg *sync.WaitGroup, topics []string, increases []string, decreases []string) error

// TopicWatcher used to watch the specified topics changed.
type TopicWatcher struct {
	lp     *glog.Logger
	client sarama.Client

	handler TopicHandler

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	// The list of topics that currently qualified.
	topics  map[string]struct{}
	regexps map[string]*regexp.Regexp
}

// NewTopicWatcher creates new TopicWatcher.
// regexTopics eg: ["^a-.*$", "^b-$"]
func NewTopicWatcher(ctx context.Context, cfg *ClientConfig, regexTopics []string, handler TopicHandler) (*TopicWatcher, error) {
	if len(regexTopics) == 0 {
		panic("TopicWatcher: must specified at least one topic")
	}
	if handler == nil {
		panic("TopicWatcher: handler can not be nil")
	}

	lp := glog.FromContext(ctx)

	lp.Info().Msg("TopicWatcher: initializing new kafka client").String("hosts", cfg.Hosts).Fire()
	client, err := sarama.NewClient(strings.Split(cfg.Hosts, ","), cfg.convert())
	if err != nil {
		lp.Error().Error("TopicWatcher: initializes kafka client error", err).Fire()
		return nil, err
	}
	return newTopicWatcherFromClient(ctx, client, regexTopics, handler)
}

func newTopicWatcherFromClient(ctx context.Context, client sarama.Client, regexTopics []string, handler TopicHandler) (*TopicWatcher, error) {
	lp := glog.FromContext(ctx)
	cctx, cancel := context.WithCancel(ctx)
	c := &TopicWatcher{
		lp:      lp,
		client:  client,
		handler: handler,
		ctx:     cctx,
		cancel:  cancel,
		wg:      new(sync.WaitGroup),
		topics:  make(map[string]struct{}),
		regexps: make(map[string]*regexp.Regexp),
	}
	c.initRegexps(regexTopics)

	lp.Debug().Msg("TopicWatcher: successfully initialized topic watcher").Fire()
	return c, nil
}

func (c *TopicWatcher) initRegexps(regexTopics []string) {
	for _, topic := range regexTopics {
		re := regexp.MustCompile(topic)
		c.regexps[topic] = re
	}
}

// returns values:
//  - topics: The list of topics that currently qualified.
//  - increases: The list of topics added compared to the previous cycle.
//  - decreases: The list of topics reduced compared to the previous cycle
func (c *TopicWatcher) fetchTopics() (topics []string, increases []string, decreases []string, err error) {
	var availTopics []string
	availTopics, err = c.client.Topics()
	if err != nil {
		return
	}

	// Match the qualified topics.
	validTopics := make(map[string]struct{})
	for _, topic := range availTopics {
		for _, re := range c.regexps {
			if re.MatchString(topic) {
				validTopics[topic] = struct{}{}
			}
		}
	}

	// Get the increased topics.
	for topic := range validTopics {
		if _, ok := c.topics[topic]; !ok {
			c.topics[topic] = struct{}{}
			increases = append(increases, topic)
		}
	}

	// Get the decreased topics.
	for topic := range c.topics {
		if _, ok := validTopics[topic]; !ok {
			delete(c.topics, topic)
			decreases = append(decreases, topic)
		} else {
			topics = append(topics, topic)
		}
	}
	return
}

// call handler in a goroutine
func (c *TopicWatcher) chenAndRun() {
	lg := c.lp

	topics, increases, decreases, err := c.fetchTopics()
	if err != nil {
		lg.Error().Error("TopicWatcher: get the available topics error", err).Fire()
		return
	}

	// No changed, skip.
	if len(increases) == 0 && len(decreases) == 0 {
		return
	}

	lg.Debug().Msg("TopicWatcher: call handler").
		Strings("topics", topics).
		Strings("increases", increases).
		Strings("decreases", decreases).
		Fire()

	err = c.handler(c.ctx, c.wg, topics, increases, decreases)
	if err != nil {
		lg.Error().Error("TopicWatcher: handler error", err).Fire()
	} else {
		lg.Debug().Msg("TopicWatcher: handler done").Fire()
	}
}

// Watch for start watch the topics changes
func (c *TopicWatcher) Watch() (err error) {
	c.wg.Add(1)
	defer c.wg.Done()

	lg := c.lp
	lg.Debug().Msg("TopicWatcher: start monitored the topics changes").Fire()

	// process the current data.
	c.chenAndRun()

	ticker := time.NewTicker(c.client.Config().Metadata.RefreshFrequency + time.Second)
LOOP:
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				break LOOP
			}
			c.chenAndRun()
		case <-c.ctx.Done():
			break LOOP
		}
	}

	ticker.Stop()
	return
}

func (c *TopicWatcher) Close() (err error) {
	if c == nil {
		return
	}
	c.cancel()

	c.lp.Debug().Msg("TopicWatcher: wait for the watcher to close").Fire()

	err = c.client.Close()
	if err != nil && err != sarama.ErrClosedClient {
		c.lp.Error().Error("TopicWatcher: close watcher error", err).Fire()
		return
	}

	c.wg.Wait()

	c.lp.Debug().Msg("TopicWatcher: watcher successful closed").Fire()
	return
}
