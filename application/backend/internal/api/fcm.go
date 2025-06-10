package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/halftone/internal/fcm"
)

func (a *api) SubscribeToPush(c *fiber.Ctx) error {
	var req fcm.SubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing subscription request: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Token == "" || req.UserID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Token and UserID are required",
		})
	}

	// Subscribe user
	if err := a.fcmService.Subscribe(&req); err != nil {
		log.Printf("Error subscribing user %s: %v", req.UserID, err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to subscribe to push notifications",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully subscribed to push notifications",
	})
}

func (a *api) SendPushMessage(c *fiber.Ctx) error {
	var req fcm.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing send message request: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Message == nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Message is required",
		})
	}

	if req.Message.Title == "" || req.Message.Body == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Message title and body are required",
		})
	}

	if len(req.UserIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "At least one target (userIds or tokens) is required",
		})
	}

	if err := a.fcmService.SendMessage(c.Context(), &req); err != nil {
		log.Printf("Error sending push message: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to send push message",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Push message sent successfully",
	})
}
