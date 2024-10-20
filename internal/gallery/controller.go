package gallery

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/michalK00/sg-qr/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GalleryController struct {
	storage *GalleryStorage
}

func NewGalleryController(storage *GalleryStorage) *GalleryController {
	return &GalleryController{
		storage: storage,
	}
}

type createGalleryRequest struct {
	Name string `json:"name" example:"Example Gallery"`
	ExpiryDate time.Time `json:"expiry_date" example:"2023-12-31T23:59:59Z"`
}

type createGalleryResponse struct {
	ID string `json:"id"`
}

// @Summary Get all galleries of a collection.
// @Description gets all galleries of a collection with id.
// @Tags collections
// @Accept */*
// @Produce json
// @Param id path string true "Collection ID"
// @Success 200 {array} domain.GalleryDB
// @Router /collections/{id}/galleries [get]
func (c *GalleryController) getAll(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Failed to parse collection id",
		}) 
	}
	
	galleries, err := c.storage.getAllGalleries(ctx.Context(), collectionId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch galleries",
		}) 
	}
	
	return ctx.JSON(galleries)
}


// @Summary Create one gallery
// @Description Creates one gallery in collection with id
// @Tags collections
// @Accept json
// @Produce json
// @Param gallery body createGalleryRequest true "Gallery to create"
// @Param id path string true "Collection ID" example:"671442a11fd0c5eb46b5a3fa"
// @Success 201 {object} createGalleryResponse
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /collections/{id}/galleries [post]
func (c *GalleryController) createGallery(ctx *fiber.Ctx) error {
	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse collection id",
		}) 
	}

	var req createGalleryRequest
	fmt.Println(string(ctx.Body()))
	if err := ctx.BodyParser(&req); err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		}) 
	}

	id, err := c.storage.createGallery(ctx.Context(), collectionId, req.Name, req.ExpiryDate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create gallery",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(createGalleryResponse{
		ID: id,
	})
}