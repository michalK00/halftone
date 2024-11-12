package galleries

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/config"
)

func AddRoutes(app *fiber.App, controller *GalleryController, config config.EnvVars) {
	gallery := app.Group("/api/v1")

	// middlewares
	// authMiddleware := auth.NewAuthMiddleware(config)

	// routes
	gallery.Get("/collections/:collectionId/galleries", controller.getGalleries)
	gallery.Get("/collections/:collectionId/galleryCount", controller.getGalleryCount)
	gallery.Post("/collections/:collectionId/galleries", controller.createGallery)

	//gallery.Post("/:collectionId/galleries/batch")
	//gallery.Delete("/:collectionId/galleries/batch")

	gallery.Get("/galleries/:galleryId", controller.getGallery)
	gallery.Put("/galleries/:galleryId", controller.updateGallery)
	gallery.Delete("/galleries/:galleryId", controller.deleteGallery)

	gallery.Post("/galleries/:galleryId/qr", controller.generateQr)

	//TODO client routes
	//client := app.Group("/api/v1/client")
	//
	//client.Get("/galleries/:galleryId")
	//client.Post("/galleries/:galleryId")
	//client.Get("/galleries/:galleryId/photos")
	//client.Get("/galleries/:galleryId/photos/:photoId")

}
