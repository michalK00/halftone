package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectionDB struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}

type CollectionRepository interface {
	CollectionExists(ctx context.Context, name primitive.ObjectID) (bool, error)
	GetCollection(ctx context.Context, collectionId primitive.ObjectID) (CollectionDB, error)
	GetCollections(ctx context.Context) ([]CollectionDB, error)
	CreateCollection(ctx context.Context, name string) (string, error)
	DeleteCollection(ctx context.Context, collectionId primitive.ObjectID) error
	UpdateCollection(ctx context.Context, collectionId primitive.ObjectID, name string) (CollectionDB, error)
}