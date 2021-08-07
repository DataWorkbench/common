package ghttp

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// NewWsUpgrader return an instance of websocket.Upgrader with default config.
func NewWsUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		HandshakeTimeout: time.Second * 5,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
		WriteBufferPool:  &sync.Pool{},
		Subprotocols:     nil,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			http.Error(w, reason.Error(), status)
		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
}

// NewWsDialer return an instance of websocket.Dialer with default config.
func NewWsDialer() *websocket.Dialer {
	return &websocket.Dialer{
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
}
