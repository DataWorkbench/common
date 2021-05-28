package kafkawrap

import (
	"context"
	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

func Consume(ctx context.Context, client sarama.Client, topic string, msgmain chan *sarama.ConsumerMessage) {

	lp := glog.FromContext(ctx)
	lp.Info().Msg("consume kafka").Fire()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		lp.Error().Error("create consumer error", err).Fire()
	}
	defer consumer.Close()

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		lp.Error().Error("get consumer partition error", err).Fire()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Println(err)
			lp.Error().Error("consume partition data error", err).Fire()
		}
		defer pc.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				msgmain <- msg
			}
		}(pc)
	}

	wg.Wait()
}

func ConsumeWithGroup(ctx context.Context, client sarama.ConsumerGroup, topic string, msg chan *sarama.ConsumerMessage) {

	lp := glog.FromContext(ctx)
	lp.Info().Msg("consumer group up and running...").Fire()

	consumer := Consumer{
		ready: make(chan bool),
		msg:   make(chan *sarama.ConsumerMessage),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		//defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{topic}, &consumer); err != nil {
				lp.Error().Error("consumer group error", err).Fire()
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	for message := range consumer.msg {
		msg <- message
	}

	wg.Wait()
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	msg   chan *sarama.ConsumerMessage
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		consumer.msg <- message
	}

	return nil
}
