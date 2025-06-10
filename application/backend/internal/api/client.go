package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/halftone/internal/aws"
	"github.com/michalK00/halftone/internal/domain"
	"github.com/michalK00/halftone/internal/fcm"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Get gallery information (client access)
// @Description Gets gallery information for clients with valid access token
// @Tags client
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID"
// @Param Authorization header string true "Access token" example:"Bearer your-access-token"
// @Success 200 {object} domain.GalleryDB
// @Failure 401 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/client/galleries/{galleryId} [get]
func (a *api) clientGetGalleryHandler(ctx *fiber.Ctx) error {
	_, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	// Gallery is already validated and stored in context by middleware
	gallery := ctx.Locals("gallery").(domain.GalleryDB)

	// Remove sensitive information before sending to client
	gallery.UserId = ""
	gallery.Sharing.AccessToken = ""

	return ctx.JSON(gallery)
}

// @Summary Create order (client access)
// @Description Creates a new order for a gallery with valid access token
// @Tags client
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID"
// @Param Authorization header string true "Access token" example:"Bearer your-access-token"
// @Param request body createOrderRequest true "Order creation request"
// @Success 201 {object} createOrderResponse
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/client/galleries/{galleryId} [post]
func (a *api) clientCreateOrderHandler(ctx *fiber.Ctx) error {
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	//Check if some an order for this gallery exists
	exists, err := a.orderRepo.OrderExistsForGallery(ctx.Context(), galleryId)
	if err != nil {
		return ServerError(ctx, err, "failed to check order")
	}
	if exists {
		return BadRequest(ctx, fmt.Errorf("order already exists"))
	}
	var req createOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}

	// Convert photo IDs
	photoIds := make([]primitive.ObjectID, len(req.PhotoIDs))
	for i, photoIdStr := range req.PhotoIDs {
		photoId, err := primitive.ObjectIDFromHex(photoIdStr)
		if err != nil {
			return BadRequest(ctx, errors.New("invalid photo ID"))
		}
		photoIds[i] = photoId
	}

	validPhotos, err := a.photoRepo.VerifyPhotosInGallery(ctx.Context(), galleryId, photoIds)
	if err != nil {
		return ServerError(ctx, err, "Failed to verify photos")
	}
	if !validPhotos {
		return BadRequest(ctx, errors.New("some photos do not belong to this gallery"))
	}

	// Create order
	orderId, err := a.orderRepo.CreateOrder(ctx.Context(), galleryId, req.ClientEmail, req.Comment, photoIds)
	if err != nil {
		return ServerError(ctx, err, "Failed to create order")
	}

	gallery, err := a.galleryRepo.GetGalleryByID(ctx.Context(), galleryId)
	if err != nil {
		return ServerError(ctx, err, "Failed to fetch gallery")
	}

	msgReq := &fcm.SendMessageRequest{
		Message: &fcm.PushMessage{
			Title: "New order",
			Body:  fmt.Sprintf("You have 1 new order"),
		},
		UserIDs: []string{gallery.UserId},
	}

	err = a.fcmService.SendMessage(msgReq)
	if err != nil {
		fmt.Printf("Failed to send push notification: %v\n", err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(createOrderResponse{
		ID: orderId,
	})
}

type clientPhotoResponse struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	OriginalFilename string             `bson:"originalFilename" json:"originalFilename"`
	Url              string             `bson:"url" json:"url"`
	ThumbnailUrl     string             `bson:"thumbnailUrl" json:"thumbnailUrl"`
}

// @Summary Get gallery photos (client access)
// @Description Gets all photos in a gallery for clients with valid access token
// @Tags client
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID"
// @Param Authorization header string true "Access token" example:"Bearer your-access-token"
// @Success 200 {array} domain.PhotoDB
// @Failure 401 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/client/galleries/{galleryId}/photos [get]
func (a *api) clientGetGalleryPhotosHandler(ctx *fiber.Ctx) error {
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	photos, err := a.photoRepo.GetSharedPhotosByGallery(ctx.Context(), galleryId)
	if err != nil {
		return ServerError(ctx, err, "Failed to fetch photos")
	}

	clientPhotos := make([]clientPhotoResponse, len(photos))
	for index, photo := range photos {
		url, err := aws.GetObjectUrl(photo.ObjectKey)
		if err != nil {
			return ServerError(ctx, err, "Failed to get url")
		}
		thumbnailUrl, err := aws.GetObjectUrl(photo.ThumbnailObjectKey)
		if err != nil {
			thumbnailUrl = ""
		}
		clientPhotos[index] = clientPhotoResponse{
			ID:               photo.ID,
			OriginalFilename: photo.OriginalFilename,
			Url:              url,
			ThumbnailUrl:     thumbnailUrl,
		}
	}

	return ctx.JSON(clientPhotos)
}

// @Summary Get specific photo (client access)
// @Description Gets a specific photo from a gallery for clients with valid access token
// @Tags client
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID"
// @Param photoId path string true "Photo ID"
// @Param Authorization header string true "Access token" example:"Bearer your-access-token"
// @Success 200 {object} domain.PhotoDB
// @Failure 401 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/client/galleries/{galleryId}/photos/{photoId} [get]
func (a *api) clientGetPhotoHandler(ctx *fiber.Ctx) error {
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	photoId, err := primitive.ObjectIDFromHex(ctx.Params("photoId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	// Get photo and verify it belongs to the gallery
	photo, err := a.photoRepo.GetSharedPhotoById(ctx.Context(), photoId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFound(ctx, err)
		}
		return ServerError(ctx, err, "Failed to fetch photo")
	}

	// Verify photo belongs to the gallery
	if photo.GalleryId != galleryId {
		return NotFound(ctx, errors.New("photo not found in this gallery"))
	}

	url, err := aws.GetObjectUrl(photo.ClientObjectKey)
	if err != nil {
		return ServerError(ctx, err, "Failed to get url")
	}
	thumbnailUrl, err := aws.GetObjectUrl(photo.ThumbnailObjectKey)
	if err != nil {
		return ServerError(ctx, err, "Failed to get url")
	}
	clientPhoto := clientPhotoResponse{
		ID:               photo.ID,
		OriginalFilename: photo.OriginalFilename,
		Url:              url,
		ThumbnailUrl:     thumbnailUrl,
	}

	return ctx.JSON(clientPhoto)
}
