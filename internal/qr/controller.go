package qr

import (
	"github.com/gofiber/fiber/v2"
)

type QrController struct {
	service *QrService
}

func NewQrController(service *QrService) *QrController {
	return &QrController{
		service: service,
	}
}

// @Description Request body for generating a QR code
type QrGenerationRequest struct {
	Url string `json:"url" example:"https://example.com"` // URL to be encoded in the QR code
}

const qrSize int = 256

// @Summary Generate QR code
// @Description Generate a QR code from a given URL
// @Tags QR
// @Accept json
// @Produce json
// @Param request body QrGenerationRequest true "QR Generation Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /qr [post]
func (c *QrController) generate(ctx *fiber.Ctx) error {
	var req QrGenerationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to generate qr code",
		})
	}

	_, err := c.service.generateQr(simpleQrCode{Content: req.Url, Size: qrSize})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate qr code",
		})
	}
	// log.Print(code)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "GG",
	})
}
