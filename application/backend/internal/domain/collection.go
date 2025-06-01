package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CollectionDB struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	UserId    string             `bson:"userId" json:"userId"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CollectionRepository interface {
	CollectionExists(ctx context.Context, name primitive.ObjectID, userId string) (bool, error)
	GetCollection(ctx context.Context, collectionId primitive.ObjectID, userId string) (CollectionDB, error)
	GetCollections(ctx context.Context, userId string) ([]CollectionDB, error)
	CreateCollection(ctx context.Context, name, userId string) (string, error)
	DeleteCollection(ctx context.Context, collectionId primitive.ObjectID, userId string) error
	UpdateCollection(ctx context.Context, collectionId primitive.ObjectID, name string, userId string) (CollectionDB, error)
}
