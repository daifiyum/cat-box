package services

import (
	U "github.com/daifiyum/cat-box/common"
	"github.com/gofiber/contrib/websocket"
)

func ShowBoxLogs(c *websocket.Conn) {
	logChan := U.Broadcaster.Subscribe(100)
	defer func() {
		U.Broadcaster.Unsubscribe(logChan)
		c.Close()
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case logLine, ok := <-logChan:
			if !ok {
				return
			}
			if err := c.WriteMessage(websocket.TextMessage, []byte(logLine)); err != nil {
				return
			}
		}
	}
}
