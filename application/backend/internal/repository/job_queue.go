package repository

import (
	"context"
	"fmt"
	"github.com/michalK00/halftone/internal/domain"
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

	jobKey := fmt.Sprintf("job:%s", job.ID.Hex())
	jobMap := map[string]interface{}{
		"id":           job.ID.Hex(),
		"type":         job.Type,
		"queue":        job.Queue,
		"status":       string(job.Status),
		"payload":      string(job.Payload),
		"created_at":   job.CreatedAt,
		"scheduled_at": job.ScheduledAt,
		"retries":      job.Retries,
	}
	if job.StartedAt != nil {
		jobMap["started_at"] = job.StartedAt
	}
	if job.CompletedAt != nil {
		jobMap["completed_at"] = job.CompletedAt
	}
	if !job.WorkerID.IsZero() {
		jobMap["worker_id"] = job.WorkerID.Hex()
	}
	if job.Error != "" {
		jobMap["error"] = job.Error
	}
	pipe.HSet(ctx, jobKey, jobMap)

	pipe.LPush(ctx, queueKey(job.Queue), jobKey)

	_, err := pipe.Exec(ctx)
	return err
}

// TODO implement
func (s *RedisJobQueue) PullJob(ctx context.Context, queueType string) (*domain.Job, error) {
	return nil, nil
}
