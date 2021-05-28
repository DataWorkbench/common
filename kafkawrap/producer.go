package kafkawrap

import (
	"context"
	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
	"sync"
)

func SyncProduce(ctx context.Context, client sarama.Client, topic string, msg string) (pid int32, offset int64, err error) {

	lp := glog.FromContext(ctx)
	lp.Info().Msg("kafka produce").Fire()
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		lp.Error().Error("create sync producer error", err).Fire()
	}
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	pid, offset, err = producer.SendMessage(message)

	return
}

func AsyncProduce(ctx context.Context, client sarama.Client, topic string, msg string) {

	lp := glog.FromContext(ctx)
	lp.Info().Msg("kafka async producer ").Fire()

	producer, err := sarama.NewAsyncProducerFromClient(client)
	//defer producer.AsyncClose()
	if err != nil {
		lp.Error().Error("create async producer error", err).Fire()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	go func(sarama.AsyncProducer) {
		producer.Input() <- message
	}(producer)

	wg.Wait()

}
