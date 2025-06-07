package aws

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type Config struct {
	Region          string
	Config          aws.Config
	CredentialScope string
}

type ServiceClient interface {
	Initialize(ctx context.Context, cfg *Config) error
}

type Client struct {
	Config  *Config
	S3      *S3Client
	Cognito *CognitoClient
}

var (
	client    *Client
	clientErr error
	once      sync.Once
)

func GetClient() (*Client, error) {
	once.Do(func() {
		client = &Client{}
		clientErr = initClient(context.Background(), client)
	})

	return client, clientErr
}

func initClient(ctx context.Context, client *Client) error {
	config, err := loadAWSConfig(ctx)
	if err != nil {
		return err
	}

	client.Config = config

	client.S3 = NewS3Client()
	if err := client.S3.Initialize(ctx, config); err != nil {
		return err
	}

	client.Cognito = NewCognitoClient()
	if err := client.Cognito.Initialize(ctx, config); err != nil {
		return err
	}

	return nil
}

func loadAWSConfig(ctx context.Context) (*Config, error) {
	region := os.Getenv("AWS_REGION")
	var cfg aws.Config
	var err error

	if os.Getenv("ENV") == "cloud" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	} else {
		cfg, err = loadDevConfig(ctx, region)
	}

	if err != nil {
		return nil, err
	}

	return &Config{
		Region: region,
		Config: cfg,
	}, nil
}

func loadDevConfig(ctx context.Context, region string) (aws.Config, error) {

	provider := credentials.NewStaticCredentialsProvider(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"",
	)

	return config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(provider),
		config.WithDefaultRegion(region),
	)
}
