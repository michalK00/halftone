package aws

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	uploader      *manager.Uploader
	defaultBucket string
}

type S3Object struct {
	Key  string
	Body *[]byte
}

func NewS3Client() *S3Client {
	return &S3Client{
		defaultBucket: os.Getenv("AWS_S3_NAME"),
	}
}

func (c *S3Client) Initialize(ctx context.Context, config *Config) error {
	c.client = s3.NewFromConfig(config.Config)
	c.presignClient = s3.NewPresignClient(c.client)
	c.uploader = manager.NewUploader(c.client)
	return nil
}

func (c *S3Client) ListS3Buckets(ctx context.Context) (*s3.ListBucketsOutput, error) {
	return c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
}

func (c *S3Client) UploadObject(ctx context.Context, object *S3Object) (*manager.UploadOutput, error) {
	return c.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &c.defaultBucket,
		Key:    &object.Key,
		Body:   bytes.NewBuffer(*object.Body),
	})
}

func (c *S3Client) HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	out, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &c.defaultBucket,
		Key:    &key,
	})
	return out, err
}

func (c *S3Client) DeleteObject(ctx context.Context, objectKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &c.defaultBucket,
		Key:    &objectKey,
	}

	_, err := c.client.DeleteObject(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func (c *S3Client) GetObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_NAME")),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			os.Getenv("AWS_S3_NAME"), objectKey, err)
	}
	return request, err
}

func (c *S3Client) PostObjectRequest(ctx context.Context, objectKey string, lifetimeSecs int64, conditions []interface{}) (*s3.PresignedPostRequest, error) {

	request, err := c.presignClient.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.defaultBucket),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignPostOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
		opts.Conditions = conditions
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to postaws sdk %v:%v. Here's why: %v\n",
			os.Getenv("AWS_S3_NAME"), objectKey, err)
	}
	return request, err
}
