package gws

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/gorilla/websocket"
	"github.com/opentracing/opentracing-go"

	"github.com/DataWorkbench/common/gtrace"
	"github.com/DataWorkbench/common/qerror"
)

type Upgrader struct {
	raw    *websocket.Upgrader
	tracer opentracing.Tracer
}

// NewUpgrader return an instance of websocket.Upgrader with default config.
func NewUpgrader(ctx context.Context) *Upgrader {
	raw := &websocket.Upgrader{
		HandshakeTimeout: time.Second * 5,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
		WriteBufferPool:  &sync.Pool{},
		Subprotocols:     nil,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			//http.Error(w, reason.Error(), status)
			resp := qerror.Response{
				Code:      "WebSocketError",
				EnUS:      reason.Error(),
				ZhCN:      reason.Error(),
				Status:    status,
				RequestID: w.Header().Get(gtrace.HeaderKey),
			}
			b, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}

			w.Header().Set("Sec-Websocket-Version", "13")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(status)
			_, _ = w.Write(b)
		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}

	ws := &Upgrader{
		raw:    raw,
		tracer: gtrace.TracerFromContext(ctx),
	}
	return ws
}

func (ws *Upgrader) Upgrade(ctx context.Context, w http.ResponseWriter, r *http.Request, responseHeader http.Header) (conn *websocket.Conn, err error) {
	lg := glog.FromContext(ctx)

	//if responseHeader == nil {
	//	responseHeader = make(http.Header, 2)
	//	if tid := gtrace.IdFromContext(ctx); tid != "" {
	//		responseHeader.Set(gtrace.HeaderKey, tid)
	//	}
	//	if span := opentracing.SpanFromContext(ctx); span != nil {
	//		err = ws.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(responseHeader))
	//		if err != nil {
	//			lg.Error().Error("inject span to websocket response header error", err).Fire()
	//		}
	//	}
	//}

	lg.Debug().Msg("upgrading request to websocket").Fire()
	conn, err = ws.raw.Upgrade(w, r, responseHeader)
	if err != nil {
		lg.Error().Error("upgraded websocket error", err).Fire()
		return
	}
	lg.Debug().Msg("successful upgraded to websocket").Fire()
	return
}
