package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

type shareGalleryRequest struct {
	SharingStart  time.Time `json:"sharingStart"`
	SharingExpiry time.Time `json:"sharingExpiry"`
}

type shareGalleryResponse struct {
	GalleryId     string    `json:"galleryId"`
	AccessToken   string    `json:"accessToken"`
	ShareUrl      string    `json:"shareUrl"`
	SharingExpiry time.Time `json:"sharingExpiry"`
}

// @Summary Share Gallery
// @Description Create a shareable link for a gallery with an expiration date
// @Tags Gallery
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID" example:"671442a11fd0c5eb46b5a3fa"
// @Param request body shareGalleryRequest true "Share Gallery Request"
// @Success 200 {object} shareGalleryResponse
// @Failure 400 {object} map[string]string "Invalid request body or expiry date"
// @Failure 404 {object} map[string]string "Gallery not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/galleries/{galleryId}/share [post]
func (a *api) shareGalleryHandler(ctx *fiber.Ctx) error {
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	var req shareGalleryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}
	if req.SharingStart.After(req.SharingExpiry) {
		return BadRequest(ctx, fmt.Errorf("sharing start time is after expiry"))
	}
	if !validateSharingExpiryDate(req.SharingExpiry) {
		return BadRequest(ctx, fmt.Errorf("sharing expiry date invalid"))
	}

	accessToken, err := domain.GenerateAccessToken()

	shareJob, err := domain.NewGalleryShareJob(domain.GallerySharePayload{GalleryId: galleryId}, req.SharingStart)
	if err != nil {
		return ServerError(ctx, err, "Failed to create job")
	}
	cleanupJob, err := domain.NewGalleryCleanupJob(domain.GalleryCleanupPayload{GalleryId: galleryId}, req.SharingExpiry)
	if err != nil {
		return ServerError(ctx, err, "Failed to create job")
	}

	_, err = a.galleryRepo.UpdateGallery(ctx.Context(), galleryId,
		domain.WithSharingEnabled(true),
		domain.WithValidatedSharingExpiryDate(primitive.NewDateTimeFromTime(req.SharingExpiry)),
		domain.WithAccessToken(accessToken),
	)
	if err != nil {
		return ServerError(ctx, err, "Failed to update gallery")
	}

	_, err = a.jobRepo.CreateJob(ctx.Context(), shareJob)
	if err != nil {
		return ServerError(ctx, err, "Failed to create share job")
	}
	_, err = a.jobRepo.CreateJob(ctx.Context(), cleanupJob)
	if err != nil {
		return ServerError(ctx, err, "Failed to create cleanup job")
	}

	return ctx.Status(fiber.StatusOK).JSON(shareGalleryResponse{
		GalleryId:     galleryId.Hex(),
		AccessToken:   accessToken,
		ShareUrl:      fmt.Sprintf("%s/galleries/%s?token=%s", os.Getenv("FRONTEND_ORIGIN"), galleryId.Hex(), accessToken),
		SharingExpiry: req.SharingExpiry,
	})
}

func validateSharingExpiryDate(expiryDate time.Time) bool {
	return !expiryDate.IsZero() && !expiryDate.Before(time.Now())
}
