package api

import (
	"github.com/gofiber/fiber/v2"
	awsClient "github.com/michalK00/halftone/platform/cloud/aws"
)

func SignUp(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	client, err := awsClient.GetClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get aws client"})
	}

	_, err = client.Cognito.SignUp(c.Context(), input.Username, input.Email, input.Password, map[string]string{})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully. Please check your email for verification.",
	})
}

func SignIn(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}
	client, err := awsClient.GetClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get aws client",
		})
	}

	resp, err := client.Cognito.InitiateAuth(c.Context(), input.Username, input.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id_token":      *resp.AuthenticationResult.IdToken,
		"access_token":  *resp.AuthenticationResult.AccessToken,
		"refresh_token": *resp.AuthenticationResult.RefreshToken,
		"expires_in":    resp.AuthenticationResult.ExpiresIn,
	})
}

func VerifyAccount(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	client, err := awsClient.GetClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get aws client",
		})
	}
	_, err = client.Cognito.ConfirmSignUp(c.Context(), input.Username, input.Code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account verified successfully",
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	client, err := awsClient.GetClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get aws client",
		})
	}

	resp, err := client.Cognito.RefreshToken(c.Context(), input.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id_token":     *resp.AuthenticationResult.IdToken,
		"access_token": *resp.AuthenticationResult.AccessToken,
		"expires_in":   resp.AuthenticationResult.ExpiresIn,
	})
}
