package collection

import (
	"context"

	"github.com/michalK00/sg-qr/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionStorage struct {
	db *mongo.Database
}

func NewCollectionStorage(db *mongo.Database) *CollectionStorage {
	return &CollectionStorage{
		db: db,
	}
}

func (s* CollectionStorage) createCollection(ctx context.Context, name string) (string, error) {
	collection := s.db.Collection("collections")

	result, err := collection.InsertOne(ctx, bson.M{"name": name, "galleries": []domain.GalleryDB{}})
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
} 

func (s* CollectionStorage) getAllCollections(ctx context.Context) ([]domain.CollectionDB, error) {
	collection := s.db.Collection("collections")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	collections := make([]domain.CollectionDB, 0)
	if err = cursor.All(ctx, &collections); err != nil {
		return nil, err
	}
	
	return collections, nil
}