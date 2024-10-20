package qr

import (
	"log"

	"github.com/skip2/go-qrcode"
)

type QrService struct {}

type simpleQrCode struct {
	Content string
	Size int
}

func (s *QrService) generateQr(qrParams simpleQrCode) ([]byte, error) {
	png, err := qrcode.Encode(qrParams.Content, qrcode.Medium, qrParams.Size)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Testing --
	err = qrcode.WriteFile(qrParams.Content, qrcode.Medium, qrParams.Size, "qr.png")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// --
	return png, nil
}