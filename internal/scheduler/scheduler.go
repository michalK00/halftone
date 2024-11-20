package scheduler

import (
	"context"
	"fmt"
	"github.com/michalK00/sg-qr/internal/domain"
	"github.com/redis/go-redis/v9"
)

func queueKey(queueType string) string {
	return fmt.Sprintf("queue:%s", queueType)
}

func PushJob(ctx context.Context, rdb *redis.Client, job domain.Job) error {
	pipe := rdb.Pipeline()

	jobKey := fmt.Sprintf("job:%s", job.ID)
	pipe.HSet(ctx, jobKey, job)

	rdb.LPush(ctx, queueKey(job.Queue), job)

	_, err := pipe.Exec(ctx)
	return err
}

// TODO implement
func PullJob(ctx context.Context, rdb *redis.Client, queueType string) (*domain.Job, error) {
	return nil, nil
}
