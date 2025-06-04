package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func FiberMiddleware(a *fiber.App) {
	allowOrigins := os.Getenv("CLIENT_ORIGIN")
	if allowOrigins == "" {
		allowOrigins = "*"
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
