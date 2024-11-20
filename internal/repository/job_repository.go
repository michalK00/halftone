package repository

import (
	"context"
	"github.com/michalK00/sg-qr/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoJob struct {
	db *mongo.Database
}

func NewMongoJob(db *mongo.Database) *MongoJob {
	return &MongoJob{
		db: db,
	}
}

func (s *MongoJob) GetJobsDue(ctx context.Context) ([]domain.Job, error) {
	collection := s.db.Collection("jobs")
	filter := bson.D{
		{"scheduled_at", bson.D{{"$lte", time.Now()}}},
	}

	var result []domain.Job
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MongoJob) CreateJob(ctx context.Context, job *domain.Job) (string, error) {
	collection := s.db.Collection("jobs")
	jobID := primitive.NewObjectID()

	job.ID = jobID

	_, err := collection.InsertOne(ctx, job)
	if err != nil {
		return "", err
	}
	return jobID.Hex(), nil
}

func (s *MongoJob) RescheduleJob(ctx context.Context, jobId primitive.ObjectID, updatedScheduledAt time.Time) (domain.Job, error) {
	collection := s.db.Collection("jobs")

	filter := bson.D{{"_id", jobId}}
	update := bson.D{
		{"$set", bson.D{
			{"scheduled_at", updatedScheduledAt},
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var job domain.Job
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&job)
	return job, err
}

func (s *MongoJob) DeleteJob(ctx context.Context, jobId primitive.ObjectID) error {
	coll := s.db.Collection("jobs")
	_, err := coll.DeleteOne(ctx, bson.D{{"_id", jobId}})
	return err
}
