package collections

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// @Summary Create collection
// @Description Creates one collection
// @Tags collections
// @Accept json
// @Produce json
// @Param collections body createCollectionRequest true "Collection to create"
// @Success 201 {object} createCollectionResponse
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/collections [post]
func (c *CollectionController) createCollection(ctx *fiber.Ctx) error {

	var req createCollectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return util.BadRequest(ctx, err)
	}

	id, err := c.storage.createCollection(ctx.Context(), req.Name)
	if err != nil {
		return util.ServerError(ctx, err, "Failed to create collection")
	}

	return ctx.Status(fiber.StatusCreated).JSON(createCollectionResponse{
		ID: id,
	})
}

// @Summary Get collections
// @Description Gets all collections
// @Tags collections
// @Accept */*
// @Produce json
// @Success 200 {array} domain.CollectionDB
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/collections [get]
func (c *CollectionController) getCollections(ctx *fiber.Ctx) error {

	collections, err := c.storage.getCollections(ctx.Context())
	if err != nil {
		return util.ServerError(ctx, err, "Failed to get collections")
	}

	return ctx.JSON(collections)
}

// @Summary Get collection
// @Description Gets specific collection
// @Tags collections
// @Accept */*
// @Produce json
// @Param collectionId path string true "Collection ID" example:"671442a11fd0c5eb46b5a3fa"
// @Success 200 {array} domain.CollectionDB
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/collections/{collectionId} [get]
func (c *CollectionController) getCollection(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return util.NotFound(ctx, err)
	}

	collection, err := c.storage.getCollection(ctx.Context(), collectionId)
	if err != nil {
		return util.ServerError(ctx, err, "Failed to get collection")
	}

	return ctx.JSON(collection)
}

type updateCollectionRequest struct {
	Name string `json:"name"`
}

// @Summary Update collection
// @Description Updates specific collection
// @Tags collections
// @Accept */*
// @Produce json
// @Param request body updateCollectionRequest true "Collection update request"
// @Param collectionId path string true "Collection ID" example:"671442a11fd0c5eb46b5a3fa"
// @Success 200 {array} domain.CollectionDB
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/collections/{collectionId} [put]
func (c *CollectionController) updateCollection(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return util.NotFound(ctx, err)
	}

	var req updateCollectionRequest
	err = ctx.BodyParser(&req)
	if err != nil {
		return util.BadRequest(ctx, err)
	}

	collection, err := c.storage.updateCollection(ctx.Context(), collectionId, req.Name)
	if err != nil {
		return util.ServerError(ctx, err, "Failed to update collection")
	}

	return ctx.Status(fiber.StatusOK).JSON(collection)
}

// @Summary Delete collection
// @Description Deletes specific collection
// @Tags collections
// @Accept */*
// @Produce json
// @Param collectionId path string true "Collection ID" example:"671442a11fd0c5eb46b5a3fa"
// @Success 204
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/collections/{collectionId} [delete]
func (c *CollectionController) deleteCollection(ctx *fiber.Ctx) error {

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusNoContent)
	}

	err = c.storage.deleteCollection(ctx.Context(), collectionId)
	if err != nil {
		return util.ServerError(ctx, err, "Failed to delete collection")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
