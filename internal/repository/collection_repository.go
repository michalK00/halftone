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

type MongoCollection struct {
	db *mongo.Database
}

func NewMongoCollection(db *mongo.Database) *MongoCollection {
	return &MongoCollection{
		db: db,
	}
}

func (s *MongoCollection) CollectionExists(ctx context.Context, collectionId primitive.ObjectID) (bool, error) {
	coll := s.db.Collection("collections")

	count, err := coll.CountDocuments(ctx, bson.D{{"_id", collectionId}}, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *MongoCollection) GetCollection(ctx context.Context, collectionId primitive.ObjectID) (domain.CollectionDB, error) {
	coll := s.db.Collection("collections")

	var collection domain.CollectionDB
	err := coll.FindOne(ctx, bson.D{{"_id", collectionId}}).Decode(&collection)
	if err != nil {
		return domain.CollectionDB{}, err
	}

	return collection, nil
}

func (s *MongoCollection) GetCollections(ctx context.Context) ([]domain.CollectionDB, error) {
	collection := s.db.Collection("collections")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	collections := make([]domain.CollectionDB, 0)
	if err = cursor.All(ctx, &collections); err != nil {
		return nil, err
	}

	return collections, nil
}

func (s *MongoCollection) CreateCollection(ctx context.Context, name string) (string, error) {
	collection := s.db.Collection("collections")

	col := bson.D{
		{"name", name},
		{"createdAt", primitive.NewDateTimeFromTime(time.Now())},
		{"updatedAt", primitive.NewDateTimeFromTime(time.Now())},
	}
	result, err := collection.InsertOne(ctx, col)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *MongoCollection) DeleteCollection(ctx context.Context, collectionId primitive.ObjectID) error {
	collection := s.db.Collection("collections")
	_, err := collection.DeleteOne(ctx, bson.D{{"_id", collectionId}})
	return err
}

func (s *MongoCollection) UpdateCollection(ctx context.Context, collectionId primitive.ObjectID, name string) (domain.CollectionDB, error) {
	coll := s.db.Collection("collections")

	filter := bson.D{{"_id", collectionId}}
	update := bson.D{
		{"$set", bson.D{
			{"name", name},
		}},
		{"$currentDate", bson.D{
			{"updatedAt", true},
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var collection domain.CollectionDB
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&collection)
	if err != nil {
		return domain.CollectionDB{}, err
	}
	return collection, nil

}