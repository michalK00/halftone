package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/halftone/internal/domain"
	"github.com/michalK00/halftone/internal/fcm"
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
// @Router /api/v1/galleries/{galleryId}/sharing/share [post]
func (a *api) shareGalleryHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
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

	gallery, err := a.galleryRepo.GetGallery(ctx.Context(), galleryId, userId)
	if err != nil {
		return NotFound(ctx, err)
	}
	if !sharingExpiryDatePastDue(gallery.Sharing.SharingExpiryDate) {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"message": "Sharing already active"})
	}

	accessToken, err := domain.GenerateAccessToken()

	_, err = a.galleryRepo.UpdateGallery(ctx.Context(), galleryId, userId,
		domain.WithSharing(domain.Sharing{
			SharingEnabled:    true,
			SharingExpiryDate: req.SharingExpiry,
			AccessToken:       accessToken,
			SharingUrl:        fmt.Sprintf("%s/galleries/%s?token=%s", os.Getenv("FRONTEND_ORIGIN"), galleryId.Hex(), accessToken),
		}),
	)
	if err != nil {
		return ServerError(ctx, err, "Failed to update gallery")
	}

	msgReq := &fcm.SendMessageRequest{
		Message: &fcm.PushMessage{
			Title: "Gallery Shared",
			Body:  fmt.Sprintf("Gallery %s has been shared successfully.", gallery.Name),
		},
		UserIDs: []string{ctx.Locals("userId").(string)},
	}

	err = a.fcmService.SendMessage(msgReq)
	if err != nil {
		fmt.Printf("Failed to send push notification: %v\n", err)
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
// @Router /api/v1/galleries/{galleryId}/sharing/reschedule [put]
func (a *api) rescheduleGallerySharingHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	var req rescheduleGallerySharingRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}
	if !validateSharingExpiryDate(req.SharingExpiry) {
		return BadRequest(ctx, fmt.Errorf("sharing expiry date invalid"))
	}

	gallery, err := a.galleryRepo.GetGallery(ctx.Context(), galleryId, userId)
	if err != nil {
		return NotFound(ctx, err)
	}
	if sharingExpiryDatePastDue(gallery.Sharing.SharingExpiryDate) {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"message": "Sharing already inactive"})
	}

	_, err = a.galleryRepo.UpdateGallery(ctx.Context(), galleryId, userId, domain.WithSharing(
		domain.Sharing{
			SharingEnabled:    true,
			SharingExpiryDate: req.SharingExpiry,
			AccessToken:       gallery.Sharing.AccessToken,
			SharingUrl:        fmt.Sprintf("%s/galleries/%s?token=%s", os.Getenv("FRONTEND_ORIGIN"), galleryId.Hex(), gallery.Sharing.AccessToken),
		}))
	if err != nil {
		return ServerError(ctx, err, "Failed to update gallery")
	}

	return ctx.Status(fiber.StatusOK).JSON(shareGalleryResponse{
		GalleryId:     galleryId.Hex(),
		AccessToken:   gallery.Sharing.AccessToken,
		ShareUrl:      fmt.Sprintf("%s/galleries/%s?token=%s", os.Getenv("FRONTEND_ORIGIN"), galleryId.Hex(), gallery.Sharing.AccessToken),
		SharingExpiry: req.SharingExpiry,
	})
}

// @Summary Stop gallery sharing
// @Description Immediately stops sharing a gallery and updates sharing options
// @Tags gallery sharing
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID" format(objectId)
// @Success 200 {object} domain.GalleryDB "Gallery with updated sharing status"
// @Failure 404 {object} map[string]string "Gallery not found"
// @Failure 405 {object} map[string]string "Sharing already inactive"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/galleries/{galleryId}/sharing/stop [put]
func (a *api) stopSharingGalleryHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	gallery, err := a.galleryRepo.GetGallery(ctx.Context(), galleryId, userId)
	if err != nil {
		return NotFound(ctx, err)
	}

	if sharingExpiryDatePastDue(gallery.Sharing.SharingExpiryDate) {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"message": "Sharing already inactive"})
	}

	gallery, err = a.galleryRepo.UpdateGallery(ctx.Context(), galleryId, userId, domain.WithSharing(
		domain.Sharing{
			SharingEnabled:    false,
			SharingExpiryDate: time.Time{},
		}))
	if err != nil {
		return ServerError(ctx, err, "Failed to update gallery")
	}

	return ctx.Status(fiber.StatusOK).JSON(gallery)
}

func validateSharingExpiryDate(expiryDate time.Time) bool {
	return !expiryDate.IsZero() && !expiryDate.Before(time.Now().UTC())
}

func sharingExpiryDatePastDue(expiryDate time.Time) bool {
	now := time.Now().UTC()
	expiryUTC := expiryDate.UTC()

	expiryMidday := time.Date(expiryUTC.Year(), expiryUTC.Month(), expiryUTC.Day(), 12, 0, 0, 0, time.UTC)
	nowMidday := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)

	return expiryMidday.Before(nowMidday)
}
