package qr

import (
	"context"
	custom_aws_wrapper "github.com/michalK00/sg-qr/internal/aws"
	"github.com/michalK00/sg-qr/platform/cloud/aws"
	"log"

	"github.com/skip2/go-qrcode"
)

type QrCode struct {
	Content string
	Size    int
}

type File struct {
	Name string
	Ext  string
	Body []byte
}

func GenerateQr(qrParams QrCode) ([]byte, error) {
	body, err := qrcode.Encode(qrParams.Content, qrcode.Medium, qrParams.Size)
	if err != nil {
		log.Printf("Failed QR code encoding, %v \n", err)
		return nil, err
	}

	return body, nil
}

func UploadQr(collectionId, galleryId string, file *File) (string, error) {

	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAwsClient, %v \n", err)
		return "", err
	}

	objectKey := custom_aws_wrapper.BuildObjectKey([]string{collectionId, galleryId}, file.Name, file.Ext)

	result, err := client.UploadObject(context.Background(), &aws.S3Object{Key: objectKey, Body: &file.Body})
	if err != nil {
		log.Printf("Failed UploadObject, %v \n", err)
		return "", err
	}

	return *result.Key, nil
}
