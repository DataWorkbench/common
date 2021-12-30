package zeppelin

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	messageHandler MessageHandler
	dial           *websocket.Conn
}

func NewWebSocketClient(handler MessageHandler) *WebSocketClient {
	return &WebSocketClient{
		messageHandler: handler,
	}
}

func (client *WebSocketClient) connect(url string) error {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dial, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	client.dial = dial
	defer client.dial.Close()

	done := make(chan struct{})

	ticker := time.NewTimer(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case t := <-ticker.C:
			err = client.messageHandler.onMessage(t.String())
			if err != nil {
				return err
			}
		case <-interrupt:
			return nil
		}
		select {
		case <-done:
		case <-time.After(time.Second):
		}
		return nil
	}
}
