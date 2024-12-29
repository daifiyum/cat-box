package router

import (
	"github.com/daifiyum/cat-box/api/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// webui
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
	setting.Get("", handler.GetSetting)
	setting.Post("", handler.UpdateSetting)

	// more 放很杂的api
	more := api.Group("more")
	more.Get("/logs", handler.ShowBoxLogs)
}
