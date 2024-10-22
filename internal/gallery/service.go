package gallery

import (
	"context"
	"log"

	"github.com/michalK00/sg-qr/internal/config/aws"
	"github.com/michalK00/sg-qr/internal/utils"
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

type qrCode struct {
	Content string
	Size    int
}

type file struct {
	Name string
	Ext  string
	Body []byte
}

func (s *GalleryService) generateQr(qrParams qrCode) ([]byte, error) {
	body, err := qrcode.Encode(qrParams.Content, qrcode.Medium, qrParams.Size)
	if err != nil {
		log.Printf("Failed QR code encoding, %v \n", err)
		return nil, err
	}

	return body, nil
}

func (s *GalleryService) uploadQr(collectionId, galleryId string, file *file) (string, error) {

	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAwsClient, %v \n", err)
		return "", err
	}

	path := utils.BuildObjectKey([]string{collectionId, galleryId}, file.Name, file.Ext)

	result, err := client.UploadObject(context.Background(), &aws.S3Object{Key: path, Body: &file.Body})
	if err != nil {
		log.Printf("Failed UploadObject, %v \n", err)
		return "", err
	}

	return *result.Key, nil
}

const lifeteimSecs int64 = 60

func (s *GalleryService) getObjectUrl(key string) (string, error) {

	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return "", err
	}
	presignerClient := aws.NewPresignClient(client)
	request, err := presignerClient.GetObjectUrl(context.Background(), key, lifeteimSecs)
	if err != nil {
		log.Printf("Failed GetObjectUrl, %v \n", err)
		return "", nil
	}

	return request.URL, nil
}

func (s *GalleryService) putObjectUrl(key string) (string, error) {

	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return "", err
	}
	presignerClient := aws.NewPresignClient(client)
	request, err := presignerClient.PutObjectUrl(context.Background(), key, lifeteimSecs)
	if err != nil {
		log.Printf("Failed PutObjectUrl, %v \n", err)
		return "", nil
	}

	return request.URL, nil
}