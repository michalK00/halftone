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
	// example: "2024-12-31T23:59:59Z"
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
// @Tags gallery sharing
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID" example:"671442a11fd0c5eb46b5a3fa"
// @Param request body shareGalleryRequest true "Share Gallery Request"
// @Success 200 {object} shareGalleryResponse
// @Failure 400 {object} map[string]string "Invalid request body or expiry date"
// @Failure 404 {object} map[string]string "Gallery not found"
// @Failure 405 {object} map[string]string "Sharing already active"
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
	if !validateSharingExpiryDate(req.SharingExpiry) {
		return BadRequest(ctx, fmt.Errorf("sharing expiry date invalid"))
	}

	gallery, err := a.galleryRepo.GetGallery(ctx.Context(), galleryId)
	if err != nil {
		return NotFound(ctx, err)
	}
	if gallery.SharingOptions.SharingEnabled {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"message": "Sharing already active"})
	}

	accessToken, err := domain.GenerateAccessToken()
	shareJob, err := domain.NewGalleryShareJob(domain.GallerySharePayload{GalleryId: galleryId}, time.Now())
	if err != nil {
		return ServerError(ctx, err, "Failed to create job")
	}
	shareExpiry := time.Date(req.SharingExpiry.Year(), req.SharingExpiry.Month(), req.SharingExpiry.Day()+1, 0, 0, 1, 0, time.UTC)
	cleanupJob, err := domain.NewGalleryCleanupJob(domain.GalleryCleanupPayload{GalleryId: galleryId}, shareExpiry)
	if err != nil {
		return ServerError(ctx, err, "Failed to create job")
	}

	cleanupJobId, err := a.jobRepo.CreateJob(ctx.Context(), cleanupJob)
	if err != nil {
		return ServerError(ctx, err, "Failed to create cleanup job")
	}

	_, err = a.galleryRepo.UpdateGallery(ctx.Context(), galleryId,
		domain.WithSharingOptions(domain.SharingOptions{
			SharingEnabled:    true,
			AccessToken:       accessToken,
			SharingExpiryDate: req.SharingExpiry,
			SharingUrl:        fmt.Sprintf("%s/galleries/%s?token=%s", os.Getenv("FRONTEND_ORIGIN"), galleryId.Hex(), accessToken),
			SharingCleanupJob: cleanupJobId,
		}),
	)
	if err != nil {
		return ServerError(ctx, err, "Failed to update gallery")
	}

	err = a.jobQueue.PushJob(ctx.Context(), *shareJob)
	if err != nil {
		return ServerError(ctx, err, "Failed to create share job")
	}

	return ctx.Status(fiber.StatusOK).JSON(shareGalleryResponse{
		GalleryId:     galleryId.Hex(),
		AccessToken:   accessToken,
		ShareUrl:      fmt.Sprintf("%s/galleries/%s?token=%s", os.Getenv("FRONTEND_ORIGIN"), galleryId.Hex(), accessToken),
		SharingExpiry: req.SharingExpiry,
	})
}

type rescheduleGallerySharingRequest struct {
	// example: "2024-12-31T23:59:59Z"
	SharingExpiry time.Time `json:"sharingExpiry"`
}

// @Summary Reschedule gallery sharing expiry
// @Description Updates the expiry date for a shared gallery and disables sharing
// @Tags gallery sharing
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID" format(objectId)
// @Param request body rescheduleGallerySharingRequest true "Reschedule sharing request"
// @Success 200 {object} shareGalleryResponse "Sharing successfully rescheduled"
// @Failure 404 {object} map[string]string "Gallery not found"
// @Failure 405 {object} map[string]string "Sharing already inactive"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/galleries/{galleryId}/reschedule [put]
func (a *api) rescheduleGallerySharingHandler(ctx *fiber.Ctx) error {
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	var req rescheduleGallerySharingRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}

	gallery, err := a.galleryRepo.GetGallery(ctx.Context(), galleryId)
	if err != nil {
		return NotFound(ctx, err)
	}
	if !gallery.SharingOptions.SharingEnabled {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"message": "Sharing already inactive"})
	}

	_, err = a.jobRepo.RescheduleJob(ctx.Context(), gallery.SharingOptions.SharingCleanupJob, time.Date(req.SharingExpiry.Year(), req.SharingExpiry.Month(), req.SharingExpiry.Day()+1, 0, 0, 1, 0, time.UTC))
	if err != nil {
		return ServerError(ctx, err, "Failed to reschedule job")
	}

	_, err = a.galleryRepo.UpdateGallery(ctx.Context(), galleryId, domain.WithSharingOptions(
		domain.SharingOptions{
			SharingEnabled:    true,
			AccessToken:       gallery.SharingOptions.AccessToken,
			SharingExpiryDate: req.SharingExpiry,
			SharingUrl:        fmt.Sprintf("%s/galleries/%s?token=%s", os.Getenv("FRONTEND_ORIGIN"), galleryId.Hex(), gallery.SharingOptions.AccessToken),
			SharingCleanupJob: gallery.SharingOptions.SharingCleanupJob,
		}))
	if err != nil {
		return ServerError(ctx, err, "Failed to update gallery")
	}

	return ctx.Status(fiber.StatusOK).JSON(shareGalleryResponse{
		GalleryId:     galleryId.Hex(),
		AccessToken:   gallery.SharingOptions.AccessToken,
		ShareUrl:      fmt.Sprintf("%s/galleries/%s?token=%s", os.Getenv("FRONTEND_ORIGIN"), galleryId.Hex(), gallery.SharingOptions.AccessToken),
		SharingExpiry: req.SharingExpiry,
	})
}

func validateSharingExpiryDate(expiryDate time.Time) bool {
	return !expiryDate.IsZero() && !expiryDate.Before(time.Now())
}
