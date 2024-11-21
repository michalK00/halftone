package domain

import (
	"context"
	"encoding/json"
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
	CreatedAt   time.Time          `bson:"createdAt"`
	ScheduledAt time.Time          `bson:"scheduledAt"`
	StartedAt   *time.Time         `bson:"startedAt,omitempty"`
	CompletedAt *time.Time         `bson:"completedAt,omitempty"`
	WorkerID    primitive.ObjectID `bson:"workerId,omitempty"`
	Error       string             `bson:"error,omitempty"`
	Retries     int                `bson:"retries"`
}

type JobRepository interface {
	GetJobsDue(ctx context.Context) ([]Job, error)
	CreateJob(ctx context.Context, job *Job) (string, error)
	DeleteJob(ctx context.Context, jobId primitive.ObjectID) error
	RescheduleJob(ctx context.Context, jobId primitive.ObjectID, updatedScheduledAt time.Time) (Job, error)
}

type JobQueue interface {
	PushJob(ctx context.Context, job Job) error
	PullJob(ctx context.Context, queueType string) (*Job, error)
}

type GallerySharePayload struct {
	GalleryId primitive.ObjectID `bson:"galleryId"`
}

type GalleryCleanupPayload struct {
	GalleryId primitive.ObjectID `bson:"galleryId"`
}

func NewGalleryShareJob(payload GallerySharePayload, scheduledAt time.Time) (*Job, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		ID:          primitive.NewObjectID(),
		Type:        "share",
		Queue:       "gallery",
		Status:      JobStatusActive,
		Payload:     jsonPayload,
		CreatedAt:   time.Now(),
		ScheduledAt: scheduledAt,
		Retries:     3,
	}, nil
}

func NewGalleryCleanupJob(payload GalleryCleanupPayload, scheduledAt time.Time) (*Job, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		ID:          primitive.NewObjectID(),
		Type:        "cleanup",
		Queue:       "gallery",
		Status:      JobStatusActive,
		Payload:     jsonPayload,
		CreatedAt:   time.Now(),
		ScheduledAt: scheduledAt,
		Retries:     3,
	}, nil
}
