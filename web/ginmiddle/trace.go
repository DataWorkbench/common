package ginmiddle

import (
	"context"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"

	"github.com/DataWorkbench/common/gtrace"
	"github.com/DataWorkbench/common/utils/idgenerator"
)

var (
	// Morally a const:
	ginComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "web"}
	spanKindTag     = opentracing.Tag{Key: string(ext.SpanKind), Value: ext.SpanKindEnum("http.server")}
)

// Trace returns a middleware for trace request.
//
// The function do following operations:
//   - Start a new opentracing span.
//   - Generate a unique 16-bytes string-type "trace-id".
//   - Create a new logger object with the "trace-id".
//   - Create a new context.Context with value.
//
// Get standard library's context.Context by ginmiddle.GetStdContext(*web.Context)
//
// Value in standard library's context.Context:
//   - span:    Get it by opentracing.SpanFromContext(ctx).
//	 - traceId: Get it by gtrace.IdFromContext(ctx).
//	 - logger:  Get it by glog.FromContext(ctx).
func Trace(ctx context.Context) gin.HandlerFunc {
	idGen := idgenerator.New("")
	lp := glog.FromContext(ctx)
	tracer := gtrace.TracerFromContext(ctx)

	return func(c *gin.Context) {
		var (
			err        error
			tid        string
			parentSpan opentracing.SpanContext
		)

		start := time.Now()

		// Creates a new log object.
		nl := lp.Clone()

		// Try to extract parent span from request headers.
		parentSpan, err = tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			nl.Error().Error("extract parent span from request headers error", err).Fire()
		}
		// Start a new span for this request.
		span := tracer.StartSpan(
			ParseOpName(c.Handler()), opentracing.ChildOf(parentSpan),
			spanKindTag, ginComponentTag,
		)
		// Inherit or generate trace id.
		tid = c.Request.Header.Get(gtrace.HeaderKey)
		if tid == "" {
			// Try to use trace id as the request id.
			sp, ok := span.Context().(gtrace.SpanContext)
			if ok {
				tid = sp.TraceID().String()
			} else {
				tid, err = idGen.Take()
				if err != nil {
					nl.Error().Error("generate new trace id error", err).Fire()
				}
			}
		}
		// Insert the trace id to logger field..
		nl.ResetFields().AddString(gtrace.IdKey, tid)

		// Inject span and trace id to response header.
		err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Writer.Header()))
		if err != nil {
			nl.Error().Error("inject span to response header error", err).Fire()
		}
		c.Writer.Header().Set(gtrace.HeaderKey, tid)

		// Init a new context with span.
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		// Insert trace id to context.Context.
		ctx = gtrace.ContextWithId(ctx, tid)
		// Insert logger object to context.Context.
		ctx = glog.WithContext(ctx, nl)

		// Set standard context.Context to *web.Context.Keys
		SetStdContext(c, ctx)

		// Debug logging.
		nl.Info().String("received request", c.Request.Method+" "+c.Request.Host+c.Request.RequestURI).Fire()
		nl.Debug().String("handler", c.HandlerName()).Fire()
		nl.Debug().Any("request headers", c.Request.Header).Fire()

		// Call next handler.
		c.Next()

		status := c.Writer.Status()
		if len(c.Errors) != 0 {
			if status < 500 {
				nl.Warn().String("request handle error", c.Errors.String()).Fire()
			} else {
				nl.Error().String("completed handle error", c.Errors.String()).Fire()
			}
		}
		nl.Info().Int("completed with status", status).Millisecond("elapsed", time.Since(start)).Fire()

		span.SetTag(string(ext.HTTPMethod), c.Request.Method)
		span.SetTag(string(ext.HTTPUrl), c.Request.RequestURI)
		span.SetTag(string(ext.HTTPStatusCode), status)

		if len(c.Errors) != 0 {
			ext.Error.Set(span, true)
			span.LogFields(tracerLog.Error(c.Errors.Last()))
		}

		// Finish the space.
		span.Finish()
		// close logger to release resources
		_ = nl.Close()
	}
}
