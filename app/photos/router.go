package photos

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/config"
)

func AddRoutes(app *fiber.App, controller *PhotosController, config config.EnvVars) {
	photos := app.Group("/api/v1")

	// middlewares
	// authMiddleware := auth.NewAuthMiddleware(config)

	// routes
	//photos.Get("/galleries/:galleryId/photos")
	photos.Post("/galleries/:galleryId/photos", controller.uploadPhotos)
	//photos.Delete("/galleries/:galleryId/photos")

	photos.Get("photos/:photoId")
	photos.Delete("photos/:photoId")

}
