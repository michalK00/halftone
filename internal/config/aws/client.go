package aws

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSClient struct {
	S3Client *s3.Client
}

var (
	awsClient    *AWSClient
	awsClientErr error
	once         sync.Once
)

func GetAWSClient() (*AWSClient, error) {
	once.Do(func() {
		env, err := loadConfig()
		if err != nil {
			awsClientErr = err
			return
		}

		awsClient = &AWSClient{}
		awsClientErr = initAWSClient(context.Background(), awsClient, env)
	})

	return awsClient, awsClientErr
}

func initAWSClient(ctx context.Context, client *AWSClient, env awsVars) error {
	var cfg aws.Config
	var err error

	if os.Getenv("GO_ENV") == "production" {
		cfg, err = config.LoadDefaultConfig(ctx)
	} else {
		cfg, err = loadDevConfig(ctx, env)
	}

	if err != nil {
		return err
	}

	client.S3Client = s3.NewFromConfig(cfg)

	return nil
}

func loadDevConfig(ctx context.Context, env awsVars) (aws.Config, error) {

	provider := credentials.NewStaticCredentialsProvider(env.AWS_ACCESS_KEY_ID, env.AWS_SECRET_ACCESS_KEY, "")

	return config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(provider), config.WithDefaultRegion(env.AWS_REGION))
}
