package collection

import (
	"github.com/gofiber/fiber/v2"
)

type CollectionController struct {
	storage *CollectionStorage
}

func NewCollectionController(storage *CollectionStorage) *CollectionController {
	return &CollectionController{
		storage: storage,
	}
}

type createCollectionRequest struct {
	Name string `json:"name"`
}

type createCollectionResponse struct {
	ID string `json:"id"`
}


// @Summary Create one collection
// @Description Creates one collection
// @Tags collections
// @Accept json
// @Produce json
// @Param collection body createCollectionRequest true "Collection to create"
// @Success 201 {object} createCollectionResponse
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /collections [post]
func (c *CollectionController) createCollection(ctx *fiber.Ctx) error {

	var req createCollectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// create collection
	id, err := c.storage.createCollection(ctx.Context(), req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create collection",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(createCollectionResponse{
		ID: id,
	})
}
// @Summary Get all collections
// @Description Gets all collections
// @Tags collections
// @Accept */*
// @Produce json
// @Success 200 {array} domain.CollectionDB
// @Failure 500 {object} fiber.Map
// @Router /collections [get]
func (c *CollectionController) getAll(ctx *fiber.Ctx) error {
	
	collections, err := c.storage.getAllCollections(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch collections",
		}) 
	}
	
	return ctx.JSON(collections)
}
