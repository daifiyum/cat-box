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
	for logLine := range logChan {
		err := c.WriteMessage(websocket.TextMessage, []byte(logLine))
		if err != nil {
			return
		}
	}
}
