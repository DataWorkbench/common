package gws

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/gorilla/websocket"
	"github.com/opentracing/opentracing-go"

	"github.com/DataWorkbench/common/gtrace"
)

type Dialer struct {
	raw    *websocket.Dialer
	tracer opentracing.Tracer
}

// NewDialer return an instance of websocket.Dialer with default config.
func NewDialer(ctx context.Context) *Dialer {
	raw := &websocket.Dialer{
		NetDial:           nil,
		NetDialContext:    nil,
		Proxy:             nil,
		TLSClientConfig:   nil,
		HandshakeTimeout:  time.Second * 45,
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		WriteBufferPool:   &sync.Pool{},
		Subprotocols:      nil,
		EnableCompression: false,
		Jar:               nil,
	}
	ws := &Dialer{
		raw:    raw,
		tracer: gtrace.TracerFromContext(ctx),
	}
	return ws
}

// Send is wrapper for websocket.Dialer.DialContext. To support opentracing span.
func (ws *Dialer) DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (conn *websocket.Conn, resp *http.Response, err error) {
	lg := glog.FromContext(ctx)

	if tid := gtrace.IdFromContext(ctx); tid != "" {
		requestHeader.Set(gtrace.HeaderKey, tid)
	}
	if span := opentracing.SpanFromContext(ctx); span != nil {
		err = ws.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(requestHeader))
		if err != nil {
			lg.Error().Error("inject span to websocket request header error", err).Fire()
		}
	}

	lg.Debug().String("sending request to url", urlStr).Fire()
	conn, resp, err = ws.raw.DialContext(ctx, urlStr, requestHeader)
	if err != nil {
		lg.Error().Error("send request error", err).Fire()
		return
	}
	lg.Debug().Int("successful request with status", resp.StatusCode).Fire()
	return
}
