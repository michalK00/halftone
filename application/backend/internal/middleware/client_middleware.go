package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/halftone/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

// AuthenticateClient validates the access token for client endpoints
func AuthenticateClient(galleryRepo domain.GalleryRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Extract gallery ID from params
		galleryIdStr := ctx.Params("galleryId")
		if galleryIdStr == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Gallery ID is required",
			})
		}

		galleryId, err := primitive.ObjectIDFromHex(galleryIdStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid gallery ID",
			})
		}

		// Extract token from Authorization header
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
		}

		token := parts[1]

		// Get gallery from database
		// Note: We need to fetch by ID without userId since this is client access
		gallery, err := galleryRepo.GetGalleryByID(ctx.Context(), galleryId)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Gallery not found",
				})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch gallery",
			})
		}

		// Validate sharing is enabled
		if !gallery.Sharing.SharingEnabled {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Gallery sharing is not enabled",
			})
		}

		// Validate token matches
		if gallery.Sharing.AccessToken != token {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid access token",
			})
		}

		// Check expiry date
		if !gallery.Sharing.SharingExpiryDate.IsZero() && time.Now().UTC().After(gallery.Sharing.SharingExpiryDate) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access token has expired",
			})
		}

		// Store gallery in context for use in handlers
		ctx.Locals("gallery", gallery)

		return ctx.Next()
	}
}
