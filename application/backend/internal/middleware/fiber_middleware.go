package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func FiberMiddleware(a *fiber.App) {
	clientFrontendOrigin := os.Getenv("CLIENT_FRONTEND_ORIGIN")
	adminFrontendOrigin := os.Getenv("ADMIN_FRONTEND_ORIGIN")
	allowOrigins := ""
	if clientFrontendOrigin == "" || adminFrontendOrigin == "" {
		allowOrigins = "*"
	} else {
		allowOrigins = clientFrontendOrigin + "," + adminFrontendOrigin
	}

	a.Use(
		cors.New(cors.Config{
			AllowOrigins: allowOrigins,
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		}),
		logger.New(),
	)
}
