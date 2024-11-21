package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/qr"
)

// @Description Request body for generating a QR code
type generateQrRequest struct {
	Url string `json:"url" example:"https://example.com"` // URL to be encoded in the QR code
}

const qrSize int = 256

// @Summary Generate QR code image from URL
// @Description Generates a QR code image in PNG format from a provided URL
// @Tags QR
// @Accept json
// @Produce image/png
// @Param url query string true "URL to encode in QR code"
// @Success 200 {file} byte[] "QR code image in PNG format"
// @Failure 400 {object} map[string]string "Invalid or missing URL parameter"
// @Failure 500 {object} map[string]string "Internal server error during QR generation"
// @Router /api/v1/qr [get]
func (a *api) generateQrHandler(ctx *fiber.Ctx) error {

	url := ctx.Query("url")
	if url == "" {
		return BadRequest(ctx, fmt.Errorf("url query parameter is required"))
	}

	body, err := qr.GenerateQr(qr.QrCode{
		Content: url,
		Size:    qrSize,
	})
	if err != nil {
		return ServerError(ctx, err, "Failed to generate qr code")
	}

	return ctx.Status(fiber.StatusOK).
		Type("image/png").
		Send(body)
}
