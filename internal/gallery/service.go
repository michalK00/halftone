package gallery

import (
	"context"
	"log"

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

func (s *GalleryService) uploadQr(qrName string, file *[]byte) (string, error) {
	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAwsClient, %v \n", err)
		return "", err
	}
	result, err := client.UploadObject(context.Background(), qrName, file)
	if err != nil {
		log.Printf("Failed UploadObject, %v \n", err)
		return "", err
	}

	return *result.Key, nil
}

const lifeteimSecs int64 = 60

func (s *GalleryService) getPresignedObjectUrl(key string) (string, error) {

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
