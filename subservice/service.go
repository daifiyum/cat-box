package subservice

import (
	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/router"
	"github.com/daifiyum/cat-box/task"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SubService() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	database.InitDatabase()
	router.SetupRoutes(app)
	task.Scheduler()

	app.Listen(":3000")
}
