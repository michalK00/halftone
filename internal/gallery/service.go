package gallery

import (
	"context"
	"log"
	"path"
	"strings"

	"github.com/michalK00/sg-qr/internal/config/aws"
	"github.com/skip2/go-qrcode"
)

type GalleryService struct {
	storage *GalleryStorage
}

func NewGalleryService(storage *GalleryStorage) *GalleryService {
	return &GalleryService{
		storage: storage,
	}
}

type simpleQrCode struct {
	Content string
	Size    int
}

func (s *GalleryService) generateQr(qrParams simpleQrCode) ([]byte, error) {
	body, err := qrcode.Encode(qrParams.Content, qrcode.Medium, qrParams.Size)
	if err != nil {
		log.Printf("Failed QR code encoding, %v \n", err)
		return nil, err
	}

	return body, nil
}

func buildObjectKey(dirs []string, objectName, extension string) string {

	fullPath := path.Join(append(dirs, objectName)...)

	if extension != "" {
		fullPath += "." + strings.TrimPrefix(extension, ".")
	}

	return fullPath
}

func (s *GalleryService) uploadQr(collectionId, galleryId string, file *[]byte) (string, error) {
	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAwsClient, %v \n", err)
		return "", err
	}
	path := buildObjectKey([]string{collectionId, galleryId}, "qr", ".png")
	result, err := client.UploadObject(context.Background(), path, file)
	if err != nil {
		log.Printf("Failed UploadObject, %v \n", err)
		return "", err
	}

	return *result.Key, nil
}

const lifeteimSecs int64 = 60

func (s *GalleryService) getQrUrl(key string) (string, error) {

	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return "", err
	}
	presignerClient := aws.NewPresignClient(client)
	request, err := presignerClient.GetObject(context.Background(), key, lifeteimSecs)
	if err != nil {
		log.Printf("Failed UploadObject, %v \n", err)
		return "", nil
	}

	return request.URL, nil
}
