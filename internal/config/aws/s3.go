package aws

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct {
	PresignClient *s3.PresignClient
	env           *awsVars
}

type S3Object struct {
	Key  string
	Body *[]byte
}

const maxConcurrency int = 3

func NewPresignClient(c *AWSClient) *Presigner {
	return &Presigner{
		PresignClient: s3.NewPresignClient(c.S3Client),
		env:           c.env,
	}
}

func (c *AWSClient) ListS3Buckets(ctx context.Context) (*s3.ListBucketsOutput, error) {
	return c.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
}

func (c *AWSClient) UploadObject(ctx context.Context, object *S3Object) (*manager.UploadOutput, error) {

	uploader := manager.NewUploader(c.S3Client)

	return uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &c.env.AWS_S3_NAME,
		Key:    &object.Key,
		Body:   bytes.NewBuffer(*object.Body),
	})

}

// TODO check if no errors
func (c *AWSClient) UploadManyObjects(ctx context.Context, objects *[]S3Object) error {

	uploader := manager.NewUploader(c.S3Client)

	var wg sync.WaitGroup
	errs := make(chan error, len(*objects))
	routinesLimiter := make(chan int, maxConcurrency)

	for _, object := range *objects {
		wg.Add(1)
		routinesLimiter <- 1

		// check if it can be *S3Object
		go func(object S3Object) {

			defer func() { wg.Done(); <-routinesLimiter }()
			_, err := uploader.Upload(ctx, &s3.PutObjectInput{
				Bucket: &c.env.AWS_S3_NAME,
				Key:    &object.Key,
				Body:   bytes.NewBuffer(*object.Body),
			})

			if err != nil {
				errs <- fmt.Errorf("failed to upload object %s: %v", object.Key, err)
			}
		}(object)
	}

	wg.Wait()
	close(errs)

	var uploadErrs []string

	for err := range errs {
		if err != nil {
			uploadErrs = append(uploadErrs, err.Error())
		}
	}

	if len(uploadErrs) > 0 {
		return fmt.Errorf("failed to upload some objects: %s", strings.Join(uploadErrs, "; "))
	}

	return nil
}

func (presigner Presigner) GetObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {

	request, err := presigner.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(presigner.env.AWS_S3_NAME),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			presigner.env.AWS_S3_NAME, objectKey, err)
	}
	return request, err
}

func (presigner Presigner) PutObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (*s3.PresignedPostRequest, error) {

	request, err := presigner.PresignClient.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(presigner.env.AWS_S3_NAME),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignPostOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
		opts.Conditions = []interface{}{
			[]interface{}{"starts-with", "$Content-Type", "image/"},
			[]interface{}{"content-length-range", 1, 10485760},
		}
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to postaws sdk %v:%v. Here's why: %v\n",
			presigner.env.AWS_S3_NAME, objectKey, err)
	}
	return request, err
}
