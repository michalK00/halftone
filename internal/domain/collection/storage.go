package collection

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// this is stored in mongo
type collectionDB struct {
	ID primitive.ObjectID `bson:"_id" json:"id"` 
	Name string `bson:"name" json:"name"`
}

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

	result, err := collection.InsertOne(ctx, bson.M{"name": name})
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
} 

func (s* CollectionStorage) getAllCollections(ctx context.Context) ([]collectionDB, error) {
	collection := s.db.Collection("collections")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	collections := make([]collectionDB, 0)
	if err = cursor.All(ctx, &collections); err != nil {
		return nil, err
	}
	
	return collections, nil
}