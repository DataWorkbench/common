package kafkawrap

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/Shopify/sarama"

	"github.com/DataWorkbench/glog"
)

type ConsumerConfig struct {
	Hosts    string `json:"hosts"         yaml:"hosts"         env:"HOSTS"             validate:"required"`
	GroupID  string `json:"group_id"      yaml:"group_id"      env:"GROUP_ID"          validate:"required"`
	Assignor string `json:"assignor"      yaml:"assignor"      env:"ASSIGNOR"`
}

type ConsumerGroup struct {
	group   sarama.ConsumerGroup
	ready   chan bool
	lp      *glog.Logger
	ctx     context.Context
	cancel  func()
	tracer  opentracing.Tracer
	handler map[string]func(context.Context, *sarama.ConsumerMessage)
}

// NewConsumerGroup return a sarama kafka consumer group
func NewConsumerGroup(ctx context.Context, cfg *ConsumerConfig, td TopicType, options ...Option) (*ConsumerGroup, error) {

	opts := applyOptions(options...)
	lp := glog.FromContext(ctx)
	lp.Info().Msg("consumer group client connecting to kafka").String("hosts", cfg.Hosts).Fire()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	switch cfg.Assignor {
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	}
	client, err := sarama.NewConsumerGroup(strings.Split(cfg.Hosts, ","), cfg.GroupID, config)
	if err != nil {
		return nil, err
	}
	consumer := &ConsumerGroup{group: client, lp: lp, tracer: opts.tracer}

	consumer.registerHandler(td)

	return consumer, nil
}

func (c *ConsumerGroup) ConsumeWithGroup(topics string) {
	if c == nil {
		return
	}
	c.ready = make(chan bool)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go func() {
		for {
			// sarama.ConsumerGroup.Consume() should be called inside an infinite loop,
			// when rebalance happens, the consumer session will need to be recreated to get the new claims
			if err := c.group.Consume(c.ctx, strings.Split(topics, ","), c); err != nil {
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

	c.lp.Info().Msg("consumer session").String("Topic", claim.Topic()).Int32("partition", claim.Partition()).Fire()

	for message := range claim.Messages() {
		//c.TopicMsgHandler(context.Background(), message)
		go c.topicMsgHandler(context.Background(), message)
		session.MarkMessage(message, "")
	}
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
	if err := c.group.Close(); err != nil {
		c.lp.Error().Error("closing client error", err).Fire()
	}
	c.lp.Info().Msg("sarama consumer group stopped").Fire()
}

func (c *ConsumerGroup) registerHandler(td []TopicDesc) {
	c.handler = make(map[string]func(context.Context, *sarama.ConsumerMessage))
	for _, t := range td {
		c.handler[t.Topic] = t.Handler
	}
}

//contextWithHeaders submit Consumer span and return context that holds a reference to Kafka Consumer SpanContext.
func (c *ConsumerGroup) contextWithHeaders(ctx context.Context, headers []*sarama.RecordHeader) context.Context {

	c.lp.Info().Msg("consumer tracing").Fire()
	// Copy a new logger
	cl := c.lp.Clone()

	//Extract messages headers
	mc := &MsgHeadersCarrier{}
	for _, h := range headers {
		mc.msgHeaders = append(mc.msgHeaders, *h)
		if string(h.Key) == "rid" { //grpcwrap.ctxReqIdKey
			cl.WithFields().AddString(string(h.Key), string(h.Value))
		}
	}
	producerSpan, err := c.tracer.Extract(opentracing.TextMap, mc)
	if err != nil {
		c.lp.Error().Error("consumer interceptor extract error", err).Fire()
	}

	consumerSpan := c.tracer.StartSpan(
		"Kafka Consumer",
		opentracing.ChildOf(producerSpan),
		ext.SpanKindConsumer,
	)

	// ContextWithSpan returns a new `context.Context` that holds a reference to `span`'s SpanContext.
	ctx = opentracing.ContextWithSpan(ctx, consumerSpan)

	// Insert logger to context
	ctx = glog.WithContext(ctx, cl)

	consumerSpan.Finish()

	return ctx

}

func (c *ConsumerGroup) topicMsgHandler(ctx context.Context, msg *sarama.ConsumerMessage) {
	//get context that holds a reference to Kafka Consumer SpanContext.
	ctx = c.contextWithHeaders(ctx, msg.Headers)
	//get topic handler func according to sarama.ConsumerMessage.Topic
	f := c.handler[msg.Topic]
	f(ctx, msg)
}

type messageHandler func(ctx context.Context, msg *sarama.ConsumerMessage)

type TopicType []TopicDesc

type TopicDesc struct {
	Topic   string
	Handler messageHandler
}

//消费者使用示例：(待删除)
//1:消费者应用初始化时注册了每个topic的处理函数的映射关系：
//consumer, err = kafkawrap.NewConsumerGroup(ctx, cfg.Kafka, handler.TopicHandler, kafkawrap.WithTracer(tracer))
//   topic的处理函数的映射关系如下：
//var TopicHandler = kafkawrap.TopicType{
//	{   Topic: "topicE",
//		Handler: TopicEHandler,
//	},
//	{   Topic: "topicB",
//		Handler:  TopicBHandler,
//	},
//}
//func TopicEHandler(ctx context.Context,msg *sarama.ConsumerMessage){
//	lp := glog.FromContext(ctx)
//	lp.Info().Msg("TopicEHandler").String("Topic", msg.Topic).String("msgValue", string(msg.Value)).Fire()
//}
//func TopicBHandler(ctx context.Context,msg *sarama.ConsumerMessage){ }

//2:启动消费程序：
//consumer.ConsumeWithGroup("topicE,topicB")
//由于sarama.ConsumerGroup.Consume()中c.joinGroupRequest(coordinator, topics)等创建消费组中消息者实例以及rebalance逻辑
//所以TopicDesc维护topic与处理函数的映射关系，sarama.ConsumerGroupHandler.ConsumeClaim()方法实现再调用topic处理函数

//producer_interceptor.go line23 rid暂时从grpc Metadata取
