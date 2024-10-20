package qr

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/config"
)

func AddQrRoutes(app *fiber.App, controller *QrController, config config.EnvVars) {
	gallery := app.Group("/qr")

	// middlewares
	// authMiddleware := auth.NewAuthMiddleware(config)
	
	// routes

	gallery.Post("/", controller.generate)
}