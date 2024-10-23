package middleware

import (
	"fmt"
	"github.com/michalK00/sg-qr/internal/config"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	config config.EnvVars
}

func NewAuthMiddleware(config config.EnvVars) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

func (a *AuthMiddleware) ValidateToken(ctx *fiber.Ctx) error {
	issuerURL, err := url.Parse("https://" + a.config.AUTH0_DOMAIN + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{a.config.AUTH0_AUDIENCE},
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	authHandler := ctx.Get("Authorization")
	authHandlerParts := strings.Split(authHandler, " ")
	if len(authHandlerParts) != 2 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid authorization header",
		})
	}

	// Validating token
	tokenInfo, err := jwtValidator.ValidateToken(ctx.Context(), authHandlerParts[1])
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	fmt.Println(tokenInfo)

	return ctx.Next()
}
