package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/aws"
	"github.com/michalK00/sg-qr/internal/qr"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createGalleryRequest struct {
	Name string `json:"name" example:"Example Gallery"`
}

type createGalleryResponse struct {
	ID string `json:"id"`
}

type updateGalleryRequest struct {
	Name              string             `json:"name" example:"Example Gallery"`
	SharingEnabled    bool               `json:"sharingEnabled,omitempty" example:"false"`
	SharingExpiryDate primitive.DateTime `json:"sharingExpiryDate,omitempty" example:"2022-01-01T00:00:00"`
}

// @Summary Get all galleries of a collection.
// @Description gets all galleries of a collection with collectionId.
// @Tags collections
// @Accept */*
// @Produce json
// @Param collectionId path string true "Collection ID"
// @Success 200 {array} domain.GalleryDB
// @Router /api/v1/collections/{collectionId}/galleries [get]
func (a *api) getGalleriesHandler(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	galleries, err := a.galleryRepo.GetGalleries(ctx.Context(), collectionId)
	if err != nil {
		return ServerError(ctx, err, "Failed to fetch galleries")
	}

	return ctx.JSON(galleries)
}

// @Summary Get gallery count for a collection
// @Description Returns the total number of galleries in a specific collection
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "Collection ID (MongoDB ObjectID hex string)"
// @Success 200 {object} fiber.Map "Gallery count"
// @Failure 404 {object} fiber.Map "Collection not found or invalid ID format"
// @Failure 500 {object} fiber.Map "Server error"
// @Router /api/v1/collections/{collectionId}/galleryCount [get]
func (a *api) getGalleryCountHandler(ctx *fiber.Ctx) error {
	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	count, err := a.galleryRepo.CollectionGalleryCount(ctx.Context(), collectionId)
	if err != nil {
		return ServerError(ctx, err, "Failed to fetch gallery count")
	}
	return ctx.JSON(fiber.Map{"count": count})
}

// @Summary Create one gallery
// @Description Creates one gallery in collection with collectionId
// @Tags collections
// @Accept json
// @Produce json
// @Param galleries body createGalleryRequest true "Gallery to create"
// @Param collectionId path string true "Collection ID" example:"671442a11fd0c5eb46b5a3fa"
// @Success 201 {object} createGalleryResponse
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/collections/{collectionId}/galleries [post]
func (a *api) createGalleryHandler(ctx *fiber.Ctx) error {
	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		NotFound(ctx, err)
	}

	exists, err := a.collectionRepo.CollectionExists(ctx.Context(), collectionId)
	if err != nil {
		return ServerError(ctx, err, "Failed to check if collection exists")
	}
	if !exists {
		return NotFound(ctx, fmt.Errorf("Collection does not exist"))
	}

	var req createGalleryRequest
	if err := ctx.BodyParser(&req); err != nil {
		BadRequest(ctx, err)
	}

	id, err := a.galleryRepo.CreateGallery(ctx.Context(), collectionId, req.Name)
	if err != nil {
		return ServerError(ctx, err, "Failed to create galleries")
	}

	return ctx.Status(fiber.StatusCreated).JSON(createGalleryResponse{
		ID: id,
	})
}

// @Summary Get gallery
// @Description Retrieves a specific gallery by its ID
// @Tags galleries
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID" example:"671442a11fd0c5eb46b5a3fa"
// @Success 200 {object} domain.GalleryDB
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/galleries/{galleryId} [get]
func (a *api) getGalleryHandler(ctx *fiber.Ctx) error {
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		NotFound(ctx, err)
	}
	gallery, err := a.galleryRepo.GetGallery(ctx.Context(), galleryId)
	if err != nil {
		ServerError(ctx, err, "Failed to fetch gallery")
	}
	return ctx.JSON(gallery)
}

// @Summary Update gallery
// @Description Updates an existing gallery's information
// @Tags galleries
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID" example:"671442a11fd0c5eb46b5a3fa"
// @Param request body updateGalleryRequest true "Gallery update request"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/galleries/{galleryId} [put]
func (a *api) updateGalleryHandler(ctx *fiber.Ctx) error {
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		NotFound(ctx, err)
	}
	var req updateGalleryRequest
	if err := ctx.BodyParser(&req); err != nil {
		BadRequest(ctx, err)
	}

	gallery, err := a.galleryRepo.UpdateGallery(ctx.Context(), galleryId, req.Name, req.SharingEnabled, req.SharingExpiryDate)
	if err != nil {
		return ServerError(ctx, err, "Failed to update gallery")
	}
	return ctx.Status(fiber.StatusOK).JSON(gallery)
}

// @Summary Delete gallery
// @Description Deletes a specific gallery
// @Tags galleries
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID" example:"671442a11fd0c5eb46b5a3fa"
// @Success 200 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/galleries/{galleryId} [delete]
func (a *api) deleteGalleryHandler(ctx *fiber.Ctx) error {
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusOK)
	}
	err = a.galleryRepo.DeleteGallery(ctx.Context(), galleryId)
	if err != nil {
		return ServerError(ctx, err, "Failed to delete gallery")
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// @Description Request body for generating a QR code
type generateQrRequest struct {
	Url string `json:"url" example:"https://example.com"` // URL to be encoded in the QR code
}

const qrSize int = 256

// @Summary Generate QR code
// @Description Generate a QR code from a given URL
// @Tags QR
// @Accept json
// @Produce json
// @Param request body generateQrRequest true "QR Generation Request"
// @Param collectionId path string true "Collection ID" example:"671442a11fd0c5eb46b5a3fa"
// @Param galleryId path string true "Gallery ID" example:"671442a11fd0c5eb46b5a3fa"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/galleries/{galleryId}/qr [post]
func (a *api) generateQrHandler(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	exists, err := a.galleryRepo.GalleryExists(ctx.Context(), galleryId)
	if err != nil {
		return ServerError(ctx, err, "Server error")
	}
	if !exists {
		return NotFound(ctx, errors.New("galleries does not exist"))
	}

	var req generateQrRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}

	body, err := qr.GenerateQr(qr.QrCode{Content: req.Url, Size: qrSize})
	if err != nil {
		return ServerError(ctx, err, "Failed to generate qr code")
	}

	objectKey, err := qr.UploadQr(
		collectionId.Hex(),
		galleryId.Hex(),
		&qr.File{
			Name: "qr",
			Ext:  ".png",
			Body: body,
		},
	)

	if err != nil {
		return ServerError(ctx, err, "Failed to upload qr code")
	}
	url, err := aws.GetObjectUrl(objectKey)
	if err != nil {
		return ServerError(ctx, err, "Failed to get presinged url")
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url": url,
	})
}