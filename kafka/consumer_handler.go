package kafka

import (
	"context"
	"reflect"
	"runtime"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/DataWorkbench/common/gtrace"
	"github.com/DataWorkbench/common/utils/idgenerator"
)

var (
	_ sarama.ConsumerGroupHandler = (*consumerHandler)(nil)
)

// type helpful for caller reference.
type ConsumerMessage = sarama.ConsumerMessage

// MessageHandler callback the consumed messages, these messages are come from the same topic-partition in every calls.
//
// The function parameters:
// - ctx: The value includes traceId, glog.Logger and opentracing.Span(if tracer is enabled).
//        And you can use the `<- ctx.Done()` to monitor whether ConsumerGroup is closed.
// - messages:
//	 - The `messages` always at least one message.
//   - If `BatchMode` is false, the `messages` contains only one message.
//   - If `BatchMode` is true, consumer will try to consume as many as possible at once.
//
// If non-nil error was returns, the consumer will be block and retry1 until successful.
type MessageHandler func(ctx context.Context, messages []*ConsumerMessage) (err error)

// consumerHandler implements sarama.ConsumerGroupHandler.
type consumerHandler struct {
	lp            *glog.Logger
	handler       MessageHandler
	tracer        opentracing.Tracer
	retryInterval time.Duration
	batchMode     bool
	batchMax      int

	// Initialize inside.
	idGen       *idgenerator.IDGenerator
	interceptor handlerInterceptor
}

// newConsumerHandler creates new sarama.ConsumerGroupHandler that implements by consumerHandler.
func newConsumerHandler(ctx context.Context, handler MessageHandler, options ...Option) sarama.ConsumerGroupHandler {
	if handler == nil {
		panic("consumerHandler: MessageHandler can not be nil")
	}

	opts := applyOptions(options...)

	h := &consumerHandler{
		lp:            glog.FromContext(ctx),
		handler:       handler,
		tracer:        gtrace.TracerFromContext(ctx),
		retryInterval: opts.retryInterval,
		batchMode:     opts.batchMode,
		batchMax:      opts.batchMax,
		idGen:         idgenerator.New(""),
		interceptor:   nil,
	}

	if !h.batchMode {
		h.batchMax = 1
	}

	interceptors := []handlerInterceptor{
		h.prepareHandler,
		h.retryHandler,
		h.spanHandler,
	}

	h.interceptor = chainInterceptors(interceptors)
	return h
}

// sarama calls Setup before ConsumeClaim.
func (h *consumerHandler) Setup(sess sarama.ConsumerGroupSession) (err error) {
	return
}

// sarama calls Cleanup after ConsumeClaim.
func (h *consumerHandler) Cleanup(sess sarama.ConsumerGroupSession) (err error) {
	return
}

// saram calls it when consume start.
func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	lg := h.lp.Clone()

	//lg.WithFields().AddString("member_id", sess.MemberID())
	lg.WithFields().AddString("topic", claim.Topic())
	lg.WithFields().AddInt64("partition", int64(claim.Partition()))
	lg.WithFields().AddString("handler", runtime.FuncForPC(reflect.ValueOf(h.handler).Pointer()).Name())

	lg.Debug().Msg("consumerHandler: consume claim started").Fire()

	var pos int

	messages := make([]*sarama.ConsumerMessage, h.batchMax) // make len=cap=batchMax.

	for {
		// collects messages,
		pos, err = h.collect(sess, claim, messages)
		if err != nil {
			break
		}

		// The `messages` at least one message.
		err = h.process(sess.Context(), messages[:pos])
		if err != nil {
			break
		}

		// Mark consumer cfg offset.
		sess.MarkMessage(messages[pos-1], "")
	}

	if err != context.Canceled {
		lg.Error().Msg("consumerHandler: consume claim exited").Error("error", err).Fire()
	} else {
		lg.Debug().Msg("consumerHandler: consume claim exited with context.Canceled").Fire()
	}

	// Make sure the offset committed in kafka-server.
	sess.Commit()

	// close the logger.
	_ = lg.Close()
	return
}

// collect collects messages from `claim.Messages()` and store the message to `messages`.
//
// The func will blocking until at least one message is received or the session done.
//
// If in batchMode; the func will try to fill values to `messages` max length, but return immediately if `claim.Messages()` blocking happened,
//
// The return value 'pos' represents the valid end index in `messages`.
func (h *consumerHandler) collect(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, messages []*sarama.ConsumerMessage) (pos int, err error) {
	ctx := sess.Context()
	// Block until the receive the first message.
	select {
	case msg, ok := <-claim.Messages():
		if !ok {
			if err := ctx.Err(); err != nil {
				return -1, err
			}
			return -1, errors.New("claim.Messages chan has been closed")
		}
		messages[pos] = msg
		pos++
	case <-ctx.Done():
		return -1, ctx.Err()
	}

	if !h.batchMode {
		return
	}

	// Try to collect as many messages as possible in batch mode.
	size := len(messages)
LOOP:
	for ; pos < size; pos++ {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				if err := ctx.Err(); err != nil {
					return -1, err
				}
				return -1, errors.New("claim.Messages chan has been closed")
			}
			messages[pos] = msg
		default:
			break LOOP
		}
	}
	return
}

// process the received messages.
func (h *consumerHandler) process(ctx context.Context, messages []*sarama.ConsumerMessage) (err error) {
	if h.interceptor == nil {
		return h.handler(ctx, messages)
	}
	return h.interceptor(ctx, messages, h.handler)
}

// Handle prepare to creates new log object and get trace id and saved to `context.Context`.
func (h *consumerHandler) prepareHandler(ctx context.Context, messages []*sarama.ConsumerMessage, handler MessageHandler) (err error) {
	// create a new logger object.
	nl := h.lp.Clone()

	tid := h.getTraceId(ctx, messages)
	if tid != "" {
		ctx = gtrace.ContextWithId(ctx, tid)
		nl.WithFields().AddString("tid", tid)
	}

	ctx = glog.WithContext(ctx, nl)

	err = handler(ctx, messages)

	_ = nl.Close()
	return
}

// Retry until the messages handle successes.
func (h *consumerHandler) retryHandler(ctx context.Context, messages []*sarama.ConsumerMessage, handler MessageHandler) (err error) {
	lg := glog.FromContext(ctx)
	msg := messages[len(messages)-1]

	lg.Debug().Msg("consumerHandler: received messages").
		String("topic", msg.Topic).
		Int32("partition", msg.Partition).
		Int64("offset", msg.Offset).
		Int("num", len(messages)).
		Fire()

	err = handler(ctx, messages)
	if err != nil && err != context.Canceled {
		// Retry until callback successful.
		retries := 0
		ticker := time.NewTicker(h.retryInterval)

	LOOP:
		for {
			lg.Error().Error("consumerHandler: handle messages error", err).Int("retrying", retries).Fire()

			retries++
			select {
			case <-ticker.C:
				err = handler(ctx, messages)
				if err != nil && err != context.Canceled {
					continue LOOP
				}
				break LOOP
			case <-ctx.Done():
				err = ctx.Err()
				break LOOP
			}
		}

		ticker.Stop()
	}
	return
}

// Open a trace span. References from `message.Headers` and `context.Context`.
func (h *consumerHandler) spanHandler(ctx context.Context, messages []*sarama.ConsumerMessage, handler MessageHandler) (err error) {
	producerSpan := h.spanFromMessage(ctx, messages)
	parentSpan := h.spanFromContext(ctx)

	msg := messages[len(messages)-1]
	span := h.tracer.StartSpan(
		"ConsumeMessage",
		opentracing.ChildOf(producerSpan),
		opentracing.ChildOf(parentSpan),
		ext.SpanKindConsumer,
		traceComponentTag,
		opentracing.Tags{"topic": msg.Topic, "partition": msg.Partition, "offset": msg.Offset, "message.num": len(messages)},
	)

	// ContextWithSpan returns a new `context.Context` that holds a reference to `span`'s SpanContext.
	ctx = opentracing.ContextWithSpan(ctx, span)

	// handler.
	if err = handler(ctx, messages); err != nil {
		ext.Error.Set(span, true)
		span.LogFields(tracerLog.Error(err))
	}

	// finish span.
	span.Finish()
	return
}

func (h *consumerHandler) getTraceId(ctx context.Context, messages []*sarama.ConsumerMessage) (tid string) {
	// Only get the trace id in first message.
	msg := messages[0]

	// Get the trace id from `message.Headers` that producer passed.
	for _, mh := range msg.Headers {
		if string(mh.Key) == gtrace.IdKey {
			tid = string(mh.Value)
			return
		}
	}

	// No trace id in `message.Headers` or in `batchMode`, check whether contains in Context.
	tid = gtrace.IdFromContext(ctx)

	// Generate a new trace id if not found in any where.
	if tid == "" {
		tid, _ = h.idGen.Take()
	}
	return
}

// Returns the parent span of producer from message if not in `batchMode`.
func (h *consumerHandler) spanFromMessage(ctx context.Context, messages []*sarama.ConsumerMessage) opentracing.SpanContext {
	// Only trace the first message.
	msg := messages[0]

	carrier := &opentracingCarrier{}
	for _, mh := range msg.Headers {
		carrier.headers = append(carrier.headers, *mh)
	}

	producerSpan, err := h.tracer.Extract(opentracing.TextMap, carrier)
	if err != nil {
		lg := glog.FromContext(ctx)
		// No parent span in message headers. We can ignore this error.
		if err == opentracing.ErrSpanContextNotFound {
			lg.Debug().Msg("consumerHandler: SpanContext not found in message headers").Fire()
		} else {
			lg.Error().Error("consumerHandler: extract SpanContext from message headers error", err).Fire()
		}
	}
	return producerSpan
}

func (h *consumerHandler) spanFromContext(ctx context.Context) opentracing.SpanContext {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return span.Context()
	}
	return nil
}
