package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ---400-499---

// 400
func BadRequest(ctx *fiber.Ctx, err error) error {
	log.Println(err)
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "Bad request",
	})
}

// 404
func NotFound(ctx *fiber.Ctx, err error) error {
	log.Println(err)
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "Not found",
	})
}

// ---500-599---

// 500
func ServerError(ctx *fiber.Ctx, err error, message string) error {
	log.Println(err)
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": message,
	})
}
