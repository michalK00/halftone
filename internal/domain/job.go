package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusActive   JobStatus = "active"
	JobStatusComplete JobStatus = "complete"
	JobStatusFailed   JobStatus = "failed"
)

type Job struct {
	ID          primitive.ObjectID `bson:"_id"`
	Type        string             `bson:"type"`
	Queue       string             `bson:"queue"`
	Status      JobStatus          `bson:"status"`
	Payload     []byte             `bson:"payload"`
	CreatedAt   time.Time          `bson:"created_at"`
	ScheduledAt time.Time          `bson:"scheduled_at"`
	StartedAt   *time.Time         `bson:"started_at,omitempty"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty"`
	WorkerID    primitive.ObjectID `bson:"worker_id,omitempty"`
	Error       string             `bson:"error,omitempty"`
	Retries     int                `bson:"retries"`
}

type JobRepository interface {
	GetJobsDue(ctx context.Context) ([]Job, error)
	CreateJob(ctx context.Context, job *Job) (string, error)
	DeleteJob(ctx context.Context, jobId primitive.ObjectID) error
	RescheduleJob(ctx context.Context, jobId primitive.ObjectID, updatedScheduledAt time.Time) (Job, error)
}
