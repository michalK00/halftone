package gallery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/config"
)

func AddGalleryRoutes(app *fiber.App, controller *GalleryController, config config.EnvVars) {
	gallery := app.Group("/collections/:collectionId/galleries")

	// middlewares
	// authMiddleware := auth.NewAuthMiddleware(config)

	// routes

	gallery.Get("/", controller.getAll)
	gallery.Post("/", controller.createGallery)
	gallery.Post("/:galleryId/generate", controller.generateQr)
	gallery.Post("/upload", controller.uploadPhotos)
}
