package photos

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/config"
)

func AddRoutes(app *fiber.App, config config.EnvVars) {
	gallery := app.Group("/api/v1")

	// middlewares
	// authMiddleware := auth.NewAuthMiddleware(config)

	// routes
	gallery.Get("photos/:photoId")
	gallery.Delete("photos/:photoId")

}
