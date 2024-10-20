package collection

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/config"
)

func AddCollectionRoutes(app *fiber.App, controller *CollectionController, config config.EnvVars) {
	collections := app.Group("/collections")

	// middlewares
	// authMiddleware := auth.NewAuthMiddleware(config)
	
	// routes

	collections.Get("/", controller.getAll)
	collections.Post("/", controller.createCollection)
}