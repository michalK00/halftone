package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/halftone/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type createOrderRequest struct {
	ClientEmail string   `json:"clientEmail" validate:"required,email" example:"client@example.com"`
	Comment     string   `json:"comment" example:"Please print all photos in 10x15cm format"`
	PhotoIDs    []string `json:"photoIds" validate:"required,min=1" example:"[\"671442a11fd0c5eb46b5a3fa\"]"`
}

type createOrderResponse struct {
	ID string `json:"id"`
}

type updateOrderRequest struct {
	Status  string `json:"status,omitempty" example:"completed"`
	Comment string `json:"comment,omitempty" example:"Updated comment"`
}

// @Summary Get all orders
// @Description Gets all orders for galleries owned by the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} domain.OrderDB
// @Failure 500 {object} fiber.Map
// @Router /api/v1/orders [get]
func (a *api) getOrdersHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)

	orders, err := a.orderRepo.GetOrders(ctx.Context(), userId)
	if err != nil {
		return ServerError(ctx, err, "Failed to fetch orders")
	}

	return ctx.JSON(orders)
}

// @Summary Get order by ID
// @Description Gets a specific order if it belongs to a gallery owned by the user
// @Tags orders
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} domain.OrderDB
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/orders/{orderId} [get]
func (a *api) getOrderHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	orderId, err := primitive.ObjectIDFromHex(ctx.Params("orderId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	order, err := a.orderRepo.GetOrder(ctx.Context(), orderId, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFound(ctx, err)
		}
		return ServerError(ctx, err, "Failed to fetch order")
	}

	return ctx.JSON(order)
}

// @Summary Update order
// @Description Updates an order's status or comment
// @Tags orders
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Param request body updateOrderRequest true "Order update request"
// @Success 200 {object} domain.OrderDB
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/orders/{orderId} [put]
func (a *api) updateOrderHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	orderId, err := primitive.ObjectIDFromHex(ctx.Params("orderId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	var req updateOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}

	var updateOpts []domain.OrderUpdateOption

	if req.Status != "" {
		// Validate status
		status := domain.OrderStatus(req.Status)
		if status != domain.OrderStatusPending && status != domain.OrderStatusCompleted {
			return BadRequest(ctx, errors.New("invalid status"))
		}
		updateOpts = append(updateOpts, domain.WithOrderStatus(status))
	}

	if req.Comment != "" {
		updateOpts = append(updateOpts, domain.WithOrderComment(req.Comment))
	}

	if len(updateOpts) == 0 {
		return BadRequest(ctx, errors.New("no fields to update"))
	}

	order, err := a.orderRepo.UpdateOrder(ctx.Context(), orderId, userId, updateOpts...)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFound(ctx, err)
		}
		return ServerError(ctx, err, "Failed to update order")
	}

	return ctx.JSON(order)
}

// @Summary Delete order
// @Description Deletes an order
// @Tags orders
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /api/v1/orders/{orderId} [delete]
func (a *api) deleteOrderHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	orderId, err := primitive.ObjectIDFromHex(ctx.Params("orderId"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusNoContent)
	}

	err = a.orderRepo.DeleteOrder(ctx.Context(), orderId, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		return ServerError(ctx, err, "Failed to delete order")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
