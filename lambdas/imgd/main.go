package main

import (
	"bytes"
	"context"
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
	maxFileSize = 3 * 1024 * 1024 // 3MB in bytes
)

// s3Client is the interface for S3 operations
type s3Client interface {
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func handleRequest(ctx context.Context, s3Event events.S3Event) error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	bucket := s3Event.Records[0].S3.Bucket.Name
	key := s3Event.Records[0].S3.Object.Key
	size := s3Event.Records[0].S3.Object.Size

	if size <= maxFileSize {
		log.Printf("File %s is already under size limit", key)
		return nil
	}

	ext := strings.ToLower(filepath.Ext(key))
	if !isImageExtension(ext) {
		log.Printf("File %s is not an image", key)
		return nil
	}

	if err := processImage(ctx, client, bucket, key, size); err != nil {
		return fmt.Errorf("failed to process image %s: %v", key, err)
	}

	return nil
}

func processImage(ctx context.Context, client s3Client, bucket, key string, size int64) error {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	result, err := client.GetObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to get object: %v", err)
	}
	defer result.Body.Close()

	var img image.Image
	contentType := *result.ContentType

	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(result.Body)
	case "image/png":
		img, err = png.Decode(result.Body)
	default:
		return fmt.Errorf("unsupported image type: %s", contentType)
	}
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	width := uint(img.Bounds().Dx())
	height := uint(img.Bounds().Dy())

	var buf bytes.Buffer
	scale := calculateScaleFactor(size, contentType == "image/jpeg")
	newWidth := uint(float64(width) * scale)
	newHeight := uint(float64(height) * scale)

	resized := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	switch contentType {
	case "image/jpeg":
		err = jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 85})
	case "image/png":
		err = png.Encode(&buf, resized)
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

	log.Printf("Successfully processed image %s", key)
	return nil
}

func calculateScaleFactor(currentSize int64, isJPEG bool) float64 {
	compressionFactor := 1.0
	if isJPEG {
		compressionFactor = 0.7
	}

	scaleFactor := math.Sqrt(float64(maxFileSize) / (float64(currentSize) * compressionFactor))

	return math.Min(1.0, scaleFactor*0.9)
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
