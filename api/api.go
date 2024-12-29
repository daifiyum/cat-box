package api

import (
	"log"

	"github.com/daifiyum/cat-box/api/router"
	U "github.com/daifiyum/cat-box/config"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Run() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":" + U.Port))
}
