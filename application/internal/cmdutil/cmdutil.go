package cmdutil

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
	"time"
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

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

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
