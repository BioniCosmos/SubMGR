package main

import (
	"log"

	"github.com/bionicosmos/submgr/config"
	"github.com/bionicosmos/submgr/cron"
	"github.com/bionicosmos/submgr/handlers"
	"github.com/bionicosmos/submgr/models"
	"github.com/bionicosmos/submgr/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.Init()
	models.Init()
	cron.Init()

	app := fiber.New()
	app.Use(logger.New())
	routes.Init(app)
	handlers.SessionInit()
	log.Fatal(app.Listen(config.Conf.Listen))
}
