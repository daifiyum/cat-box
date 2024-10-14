package router

import (
	"github.com/daifiyum/cat-box/subservice/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// index
	app.Static("/", "./resources/ui/sub")

	// api
	api := app.Group("/api")

	// Subscribe
	subscribe := api.Group("/subscribe")
	subscribe.Get("/", handler.GetAllSubscribe)
	subscribe.Post("/", handler.CreateSubscribe)
	subscribe.Delete("/:id", handler.DeleteSubscribe)

	subscribe.Put("/:id/active", handler.ActiveSubscribe)
	subscribe.Put("/:id/edit", handler.EditSubscribe)
	subscribe.Put("/:id/update", handler.UpdateSubscribe)
	subscribe.Put("/order", handler.OrderSubscribe)

	// setting
	setting := api.Group("/setting")
	setting.Get("/:key", handler.GetSetting)
	setting.Post("/:key", handler.UpdateSetting)
}
