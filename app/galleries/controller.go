package galleries

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"path"
)

type GalleryController struct {
	service *GalleryService
}

func NewGalleryController(service *GalleryService) *GalleryController {
	return &GalleryController{
		service: service,
	}
}

type createGalleryRequest struct {
	Name string `json:"name" example:"Example Gallery"`
}

type createGalleryResponse struct {
	ID string `json:"id"`
}

// @Summary Get all galleries of a collection.
// @Description gets all galleries of a collection with collectionId.
// @Tags collections
// @Accept */*
// @Produce json
// @Param collectionId path string true "Collection ID"
// @Success 200 {array} domain.GalleryDB
// @Router /api/v1/collections/{collectionId}/galleries [get]
func (c *GalleryController) getGalleries(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return util.NotFound(ctx, err)
	}

	galleries, err := c.service.storage.getGalleries(ctx.Context(), collectionId)
	if err != nil {
		return util.ServerError(ctx, err, "Failed to fetch galleries")
	}

	return ctx.JSON(galleries)
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
func (c *GalleryController) createGallery(ctx *fiber.Ctx) error {
	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		util.NotFound(ctx, err)
	}

	var req createGalleryRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.BadRequest(ctx, err)
	}

	id, err := c.service.storage.createGallery(ctx.Context(), collectionId, req.Name)
	if err != nil {
		return util.ServerError(ctx, err, "Failed to create galleries")
	}

	return ctx.Status(fiber.StatusCreated).JSON(createGalleryResponse{
		ID: id,
	})
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
func (c *GalleryController) generateQr(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return util.NotFound(ctx, err)
	}
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return util.NotFound(ctx, err)
	}
	exists, err := c.service.storage.galleryExists(ctx.Context(), galleryId)
	if err != nil {
		return util.ServerError(ctx, err, "Server error")
	}
	if !exists {
		return util.NotFound(ctx, errors.New("galleries does not exist"))
	}

	var req generateQrRequest
	if err := ctx.BodyParser(&req); err != nil {
		return util.BadRequest(ctx, err)
	}

	body, err := c.service.generateQr(qrCode{Content: req.Url, Size: qrSize})
	if err != nil {
		return util.ServerError(ctx, err, "Failed to generate qr code")
	}

	objectKey, err := c.service.uploadQr(
		collectionId.Hex(),
		galleryId.Hex(),
		&file{
			Name: "qr",
			Ext:  ".png",
			Body: body,
		},
	)

	if err != nil {
		return util.ServerError(ctx, err, "Failed to upload qr code")
	}
	url, err := c.service.getObjectUrl(objectKey)
	if err != nil {
		return util.ServerError(ctx, err, "Failed to get presinged url")
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url": url,
	})
}

func (c *GalleryController) uploadPhotos(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return util.NotFound(ctx, err)
	}
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return util.NotFound(ctx, err)
	}
	exists, err := c.service.storage.galleryExists(ctx.Context(), galleryId)
	if err != nil {
		return util.ServerError(ctx, err, "Server error")
	}
	if !exists {
		return util.NotFound(ctx, errors.New("galleries does not exist"))
	}

	path := path.Join(collectionId.Hex(), galleryId.Hex(), "photos")

	url, err := c.service.putObjectUrl(path)
	if err != nil {
		return util.ServerError(ctx, err, "Failed to get presinged url")
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url": url,
	})
}
