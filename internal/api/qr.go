package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/sg-qr/internal/qr"
)

// @Description Request body for generating a QR code
type generateQrRequest struct {
	Url string `json:"url" example:"https://example.com"` // URL to be encoded in the QR code
}

const qrSize int = 256

// @Summary Generate QR code
// @Description Generate a QR code from a given URL
// @Tags QR
// @Accept json
// @Produce image/png
// @Param request body generateQrRequest true "QR Generation Request"
// @Success 200 {file} byte[] "QR code image in PNG format"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/qr [post]
func (a *api) generateQrHandler(ctx *fiber.Ctx) error {

	var req generateQrRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}

	body, err := qr.GenerateQr(qr.QrCode{Content: req.Url, Size: qrSize})
	if err != nil {
		return ServerError(ctx, err, "Failed to generate qr code")
	}

	return ctx.Status(fiber.StatusOK).Type("image/png").Send(body)
}
