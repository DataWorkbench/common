# kafka

## SyncProducer
```go
package main

import (
	"context"

	"github.com/DataWorkbench/common/kafka"
	"github.com/DataWorkbench/glog"
)

func main() {
	cfg := &kafka.ProducerConfig{
		Hosts:            "127.0.0.1:9092",
		RequiredAcks:     -1,
		PartitionerClass: "hash",
	}

	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)

	// Initializes a producer.
	producer, err := kafka.NewAsyncProducer(ctx, cfg, kafka.WithTracer(nil))
	if err != nil {
		return
	}

	// Send
	err = producer.Send(ctx, "di-3", nil, kafka.StringEncoder("Hello World!"))
	if err != nil {
		return
	}
	_ = producer.Close()
}
```

## AsyncProducer
```go
package main

import (
	"context"

	"github.com/DataWorkbench/common/kafka"
	"github.com/DataWorkbench/glog"
)

func main() {
	cfg := &kafka.ProducerConfig{
		Hosts:            "127.0.0.1:9092",
		RequiredAcks:     -1,
		PartitionerClass: "hash",
	}

	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)

	// Initializes a producer.
	producer, err := kafka.NewAsyncProducer(ctx, cfg, kafka.WithTracer(nil))
	if err != nil {
		return
	}

	// Send
	err = producer.Send(ctx, "di-3", nil, kafka.StringEncoder("Hello World!"))
	if err != nil {
		return
	}
	_ = producer.Close()
}
```

## ConsumerGroup

### Consumer process one message at a time. (Defaults)
```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataWorkbench/common/kafka"
	"github.com/DataWorkbench/glog"
)

func ConsumeHandler(ctx context.Context, messages []*kafka.ConsumerMessage) (err error) {
	lg := glog.FromContext(ctx)

	msg := messages[0]
	lg.Debug().Msg("ConsumeHandler: new messages").String("key", string(msg.Key)).String("value", string(msg.Value)).Fire()
	return
}

func main() {
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)

	cfg := &kafka.ConsumerConfig{
		Hosts:            "127.0.0.1:9092",
		RefreshFrequency: time.Second * 3,
		OffsetsInitial:   0,
		BalanceStrategy:  "sticky",
	}

	watcher, err := kafka.NewConsumerGroup(ctx, "group1", cfg, ConsumeHandler, kafka.WithTracer(nil))
	if err != nil {
		return
	}

	go watcher.Consume([]string{"space-1"})

	// handle signal
	sigGroup := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	sigChan := make(chan os.Signal, len(sigGroup))
	signal.Notify(sigChan, sigGroup...)

	blockChan := make(chan struct{})

	go func() {
		sig := <-sigChan
		lp.Info().String("receive system signal", sig.String()).Fire()
		blockChan <- struct{}{}
	}()

	<-blockChan

	_ = watcher.Close()
}
```

### Consumer process many message at a time.
```go
consumer, err := kafka.NewConsumerGroup(ctx, "group1", cfg, ConsumeHandler, kafka.WithTracer(nil), kafka.WithBatchMode(true))
```

### Consumer process the dynamic topic lists.

If sets topics with a regular expression, The consumer will monitor the kafka's topics changes, 
When qualified topics changes, The automatic re-balance will be triggers.

```go
consumer, err := kafka.NewConsumerDynamic(ctx, "group1", cfg, ConsumeHandler, kafka.WithTracer(nil))
if err != nil {
    return
}
go consumer.Consume([]string{"^space.*$", "^flow-.*"})
```

### TopicWatcher

Watch the specified topics that format by regular expression; When the list of topics changed, The `kafka.TopicHandler` will be called.

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DataWorkbench/common/kafka"
	"github.com/DataWorkbench/glog"
)

func MessageHandler(ctx context.Context, messages []*kafka.ConsumerMessage) (err error) {
	lg := glog.FromContext(ctx)
	msg := messages[0]
	lg.Debug().Msg("MessageHandler: new messages").String("key", string(msg.Key)).String("value", string(msg.Value)).Fire()
	return
}

func ConsumerWatcherHandler(cfg *kafka.ConsumerConfig, msgHandler kafka.MessageHandler, options ...kafka.Option) kafka.TopicHandler {
	return func(ctx context.Context, wg *sync.WaitGroup, _ []string, increases []string, _ []string) error {
		if len(increases) == 0 {
			return nil
		}
		for i := range increases {
			wg.Add(1)
			go func(x int) {
				defer wg.Done()

				topic := increases[x]
				groupId := topic
				consumer, err := kafka.NewConsumerGroup(ctx, groupId, cfg, msgHandler, options...)
				if err != nil {
					return
				}
				_ = consumer.Consume([]string{topic})
				_ = consumer.Close()
			}(i)
		}
		return nil
	}
}

func main() {
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)

	cfg := &kafka.ConsumerConfig{
		Hosts:            "127.0.0.1:9092",
		RefreshFrequency: time.Minute * 10,
		OffsetsInitial:   -2,
		BalanceStrategy:  "roundRobin",
	}

	cliCfg := &kafka.ClientConfig{
		Hosts:            "127.0.0.1:9092",
		RefreshFrequency: time.Second * 30,
	}

	handler := ConsumerWatcherHandler(
		cfg, MessageHandler,
		kafka.WithTracer(nil),
		kafka.WithBatchMode(true),
	)

	watcher, err := kafka.NewTopicWatcher(ctx, cliCfg, []string{"^workflow-.*$"}, handler)
	if err != nil {
		return
	}

	go watcher.Watch()

	// handle signal
	sigGroup := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	sigChan := make(chan os.Signal, len(sigGroup))
	signal.Notify(sigChan, sigGroup...)

	blockChan := make(chan struct{})

	go func() {
		sig := <-sigChan
		lp.Info().String("receive system signal", sig.String()).Fire()
		blockChan <- struct{}{}
	}()

	<-blockChan

	_ = watcher.Close()
}
```


