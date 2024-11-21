package repository

import (
	"context"
	"fmt"
	"github.com/michalK00/sg-qr/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisJobQueue struct {
	rdb *redis.Client
}

func NewRedisJob(rdb *redis.Client) *RedisJobQueue {
	return &RedisJobQueue{
		rdb: rdb,
	}
}

func queueKey(queueType string) string {
	return fmt.Sprintf("queue:%s", queueType)
}

func (s *RedisJobQueue) PushJob(ctx context.Context, job domain.Job) error {
	pipe := s.rdb.Pipeline()

	jobKey := fmt.Sprintf("job:%s", job.ID)
	pipe.HSet(ctx, jobKey, job)

	pipe.LPush(ctx, queueKey(job.Queue), job)

	_, err := pipe.Exec(ctx)
	return err
}

// TODO implement
func (s *RedisJobQueue) PullJob(ctx context.Context, queueType string) (*domain.Job, error) {
	return nil, nil
}
