package cmdutil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	// Path to the AWS CA file
	caFilePath = "/opt/global-bundle.pem"

	// Timeout operations after N seconds
	connectTimeout = 5
	queryTimeout   = 30
)

func NewLogger(service string) *zap.Logger {
	env := os.Getenv("ENV")
	logger, _ := zap.NewProduction(zap.Fields(zap.String("env", env), zap.String("service", service)))

	if env == "" || env == "development" {
		logger, _ = zap.NewDevelopment()
	}
	return logger
}

func NewMongoClient() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tlsConfig, err := getCustomTLSConfig(caFilePath)
	opts := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetTLSConfig(tlsConfig)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client.Database(os.Getenv("MONGODB_NAME")), nil
}

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URI"))
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := os.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}
