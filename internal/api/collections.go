package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
func (a *api) createCollectionHandler(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId").(string)

	var req createCollectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}

	id, err := a.collectionRepo.CreateCollection(ctx.Context(), req.Name, userId)
	if err != nil {
		return ServerError(ctx, err, "Failed to create collection")
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
func (a *api) getCollectionsHandler(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId").(string)

	collections, err := a.collectionRepo.GetCollections(ctx.Context(), userId)
	if err != nil {
		return ServerError(ctx, err, "Failed to get collections")
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
func (a *api) getCollectionHandler(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId").(string)

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	collection, err := a.collectionRepo.GetCollection(ctx.Context(), collectionId, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFound(ctx, err)
		}
		return ServerError(ctx, err, "Failed to get collection")
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
func (a *api) updateCollectionHandler(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId").(string)

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	var req updateCollectionRequest
	err = ctx.BodyParser(&req)
	if err != nil {
		return BadRequest(ctx, err)
	}

	collection, err := a.collectionRepo.UpdateCollection(ctx.Context(), collectionId, req.Name, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFound(ctx, err)
		}
		return ServerError(ctx, err, "Failed to update collection")
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
func (a *api) deleteCollectionHandler(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId").(string)

	collectionId, err := primitive.ObjectIDFromHex(ctx.Params("collectionId"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusNoContent)
	}

	err = a.collectionRepo.DeleteCollection(ctx.Context(), collectionId, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		return ServerError(ctx, err, "Failed to delete collection")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
