package util

import (
	"log"

	"github.com/skip2/go-qrcode"
)

type simpleQrCode struct {
	Content string
	Size    int
}

func GenerateQr(qrParams simpleQrCode) ([]byte, error) {
	png, err := qrcode.Encode(qrParams.Content, qrcode.Medium, qrParams.Size)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return png, nil
}
