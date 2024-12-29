package handler

import (
	"github.com/daifiyum/cat-box/api/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ShowBoxLogs(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(services.ShowBoxLogs)(c)
	}
	return fiber.ErrUpgradeRequired
}
