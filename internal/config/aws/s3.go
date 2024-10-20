package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *AWSClient) ListS3Buckets(ctx context.Context) (*s3.ListBucketsOutput, error) {
	return c.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
}
