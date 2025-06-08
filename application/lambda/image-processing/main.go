package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nfnt/resize"
)

const (
	maxClientFileSize = 3 * 1024 * 1024 // 3MB in bytes for client images
	thumbnailMaxSize  = 300             // 300px max dimension for thumbnails
	clientQuality     = 85              // JPEG quality for client images
	thumbnailQuality  = 80              // JPEG quality for thumbnails
)

type PhotoUploadPayload struct {
	GalleryId string `json:"galleryId"`
	PhotoId   string `json:"photoId"`
	ObjectKey string `json:"objectKey"`
	Bucket    string `json:"bucket"`
}

type LambdaPayload struct {
	EventType    string             `json:"eventType"`
	Payload      PhotoUploadPayload `json:"payload"`
	Metadata     map[string]string  `json:"metadata"`
	DelaySeconds int                `json:"delaySeconds"`
}

type s3Client interface {
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func handleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	for _, record := range sqsEvent.Records {
		if err := processMessage(ctx, client, record); err != nil {
			log.Printf("Failed to process message %s: %v", record.MessageId, err)
			continue
		}
	}

	return nil
}

func processMessage(ctx context.Context, client s3Client, record events.SQSMessage) error {
	var lambdaPayload LambdaPayload
	if err := json.Unmarshal([]byte(record.Body), &lambdaPayload); err != nil {
		return fmt.Errorf("failed to unmarshal SQS message: %v", err)
	}

	if lambdaPayload.EventType != "photo.uploaded" {
		log.Printf("Ignoring event type: %s", lambdaPayload.EventType)
		return nil
	}

	payload := lambdaPayload.Payload
	objectKey := payload.ObjectKey

	if !isValidPhotoPath(objectKey) {
		return fmt.Errorf("invalid object key format: %s", objectKey)
	}

	ext := strings.ToLower(filepath.Ext(objectKey))
	if !isImageExtension(ext) {
		log.Printf("File %s is not a supported image", objectKey)
		return nil
	}

	if err := processPhotoVariants(ctx, client, payload.ObjectKey, payload.Bucket); err != nil {
		return fmt.Errorf("failed to process photo variants for %s: %v", payload.ObjectKey, err)
	}

	log.Printf("Successfully processed photo variants for %s", objectKey)
	return nil
}

func processPhotoVariants(ctx context.Context, client s3Client, originalKey, bucket string) error {
	img, contentType, err := downloadImage(ctx, client, originalKey, bucket)
	if err != nil {
		return fmt.Errorf("failed to download original image: %v", err)
	}

	clientKey := generateClientPath(originalKey)
	thumbnailKey := generateThumbnailPath(originalKey)

	if err := createAndUploadImage(ctx, client, img, clientKey, contentType, bucket, ImageTypeClient); err != nil {
		return fmt.Errorf("failed to create client image: %v", err)
	}

	if err := createAndUploadImage(ctx, client, img, thumbnailKey, contentType, bucket, ImageTypeThumbnail); err != nil {
		return fmt.Errorf("failed to create thumbnail: %v", err)
	}

	return nil
}

type ImageType int

const (
	ImageTypeClient ImageType = iota
	ImageTypeThumbnail
)

func downloadImage(ctx context.Context, client s3Client, key, bucket string) (image.Image, string, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	result, err := client.GetObject(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object: %v", err)
	}
	defer result.Body.Close()

	contentType := ""
	if result.ContentType != nil {
		contentType = *result.ContentType
	}

	var img image.Image
	switch contentType {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(result.Body)
	case "image/png":
		img, err = png.Decode(result.Body)
	default:
		ext := strings.ToLower(filepath.Ext(key))
		switch ext {
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(result.Body)
			contentType = "image/jpeg"
		case ".png":
			img, err = png.Decode(result.Body)
			contentType = "image/png"
		default:
			return nil, "", fmt.Errorf("unsupported image type: %s", contentType)
		}
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %v", err)
	}

	return img, contentType, nil
}

func createAndUploadImage(ctx context.Context, client s3Client, img image.Image, key, contentType, bucket string, imageType ImageType) error {
	var resized image.Image
	var quality int

	switch imageType {
	case ImageTypeClient:
		resized = resizeForClient(img)
		quality = clientQuality
	case ImageTypeThumbnail:
		resized = resizeForThumbnail(img)
		quality = thumbnailQuality
	}

	var buf bytes.Buffer
	var err error

	switch contentType {
	case "image/jpeg", "image/jpg":
		err = jpeg.Encode(&buf, resized, &jpeg.Options{Quality: quality})
	case "image/png":
		err = png.Encode(&buf, resized)
	default:
		return fmt.Errorf("unsupported content type for encoding: %s", contentType)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: &contentType,
	})

	if err != nil {
		return fmt.Errorf("failed to upload processed image: %v", err)
	}

	log.Printf("Successfully uploaded %s", key)
	return nil
}

func resizeForClient(img image.Image) image.Image {
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())

	// Calculate scale factor to keep under maxClientFileSize
	currentPixels := width * height
	maxPixels := uint(2000 * 2000) // Approximately 4MP for 3MB file size

	if currentPixels <= maxPixels {
		return img
	}

	scaleFactor := math.Sqrt(float64(maxPixels) / float64(currentPixels))
	newWidth := uint(float64(width) * scaleFactor)
	newHeight := uint(float64(height) * scaleFactor)

	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}

func resizeForThumbnail(img image.Image) image.Image {
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())

	if width <= thumbnailMaxSize && height <= thumbnailMaxSize {
		return img
	}

	var newWidth, newHeight uint
	if width > height {
		newWidth = thumbnailMaxSize
		newHeight = 0
	} else {
		newWidth = 0
		newHeight = thumbnailMaxSize
	}

	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}

func generateClientPath(originalKey string) string {
	// Convert: {collectionId}/{galleryId}/photos/{photo}
	// To: {collectionId}/{galleryId}/photos_client/{photo}
	return strings.Replace(originalKey, "/photos/", "/photos_client/", 1)
}

func generateThumbnailPath(originalKey string) string {
	// Convert: {collectionId}/{galleryId}/photos/{photo}
	// To: {collectionId}/{galleryId}/photos_client/{photo}_thumbnail
	clientPath := generateClientPath(originalKey)
	ext := filepath.Ext(clientPath)
	nameWithoutExt := strings.TrimSuffix(clientPath, ext)
	return nameWithoutExt + "_thumbnail" + ext
}

func isValidPhotoPath(key string) bool {
	// Check if the path matches the expected format: {collectionId}/{galleryId}/photos/{photo}
	parts := strings.Split(key, "/")
	return len(parts) >= 4 && parts[len(parts)-2] == "photos"
}

func isImageExtension(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

func main() {
	lambda.Start(handleRequest)
}
