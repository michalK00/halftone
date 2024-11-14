package aws

import (
	"bytes"
	"context"
	"log"
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

func (c *AWSClient) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.env.AWS_S3_NAME,
		Key:    &key,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *AWSClient) DeleteObject(ctx context.Context, objectKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &c.env.AWS_S3_NAME,
		Key:    &objectKey,
	}

	_, err := c.S3Client.DeleteObject(ctx, input)
	if err != nil {
		return err
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

func (presigner Presigner) PostObjectRequest(ctx context.Context, objectKey string, lifetimeSecs int64, conditions []interface{}) (*s3.PresignedPostRequest, error) {

	request, err := presigner.PresignClient.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(presigner.env.AWS_S3_NAME),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignPostOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
		opts.Conditions = conditions
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to postaws sdk %v:%v. Here's why: %v\n",
			presigner.env.AWS_S3_NAME, objectKey, err)
	}
	return request, err
}
