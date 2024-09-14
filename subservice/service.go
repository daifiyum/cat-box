package subservice

import (
	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/router"
	"github.com/daifiyum/cat-box/task"
	"github.com/daifiyum/cat-box/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SubService() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Output: utils.GetLogFile(),
	}))

	err := database.ConnectDB()
	if err != nil {
		utils.LogError("Failed to connect to the database")
		return
	}

	router.SetupRoutes(app)
	task.Scheduler()

	app.Listen(":3000")
}
