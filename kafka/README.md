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
watcher, err := kafka.NewConsumerGroup(ctx, "group1", cfg, ConsumeHandler, kafka.WithTracer(nil), kafka.WithBatchMode(true))
```

### Consumer monitor regex format topics.

If sets topics with a regular expression, The consumer will monitor the kafka's topics changes, 
When qualified topics changes, The automatic re-balance will be triggers.

```go
watcher, err := kafka.NewConsumerWatcher(ctx, "group1", cfg, ConsumeHandler, kafka.WithTracer(nil))
if err != nil {
    return
}
go watcher.Consume([]string{"^space.*$", "^flow-.*"})
```
