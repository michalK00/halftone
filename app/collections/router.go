package collections

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/config"
)

func AddRoutes(app *fiber.App, controller *CollectionController, config config.EnvVars) {
	route := app.Group("/api/v1")

	// middlewares
	// authMiddleware := auth.NewAuthMiddleware(config)
	// routes
	route.Get("/collections", controller.getCollections)
	route.Post("/collections", controller.createCollection)
	route.Get("/collections/:collectionId", controller.getCollection)
	route.Put("/collections/:collectionId", controller.updateCollection)
	route.Delete("/collections/:collectionId", controller.deleteCollection)

	//TODO
	//route.Get("/collections/:collectionId/qr")
}
