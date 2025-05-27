package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func FiberMiddleware(a *fiber.App) {
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "*"
	}

	a.Use(
		// Add CORS to each route
		cors.New(cors.Config{
			AllowOrigins: frontendOrigin,
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		}),
		// Add simple logger.
		logger.New(),
	)
}
