package repository

import (
	"context"
	"github.com/michalK00/halftone/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoGallery struct {
	db *mongo.Database
}

func NewMongoGallery(db *mongo.Database) *MongoGallery {
	return &MongoGallery{
		db: db,
	}
}

func (s *MongoGallery) GetGalleryByID(ctx context.Context, galleryId primitive.ObjectID) (domain.GalleryDB, error) {
	coll := s.db.Collection("galleries")

	var result domain.GalleryDB
	err := coll.FindOne(ctx, bson.M{"_id": galleryId}).Decode(&result)
	if err != nil {
		return domain.GalleryDB{}, err
	}

	return result, nil
}

func (s *MongoGallery) GalleryExists(ctx context.Context, galleryId primitive.ObjectID, userId string) (bool, error) {
	coll := s.db.Collection("galleries")

	count, err := coll.CountDocuments(ctx, bson.M{"_id": galleryId, "userId": userId}, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *MongoGallery) CollectionGalleryCount(ctx context.Context, collectionId primitive.ObjectID, userId string) (int64, error) {
	coll := s.db.Collection("galleries")
	count, err := coll.CountDocuments(ctx, bson.M{"collectionId": collectionId, "userId": userId})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *MongoGallery) GetGalleries(ctx context.Context, collectionId primitive.ObjectID, userId string) ([]domain.GalleryDB, error) {
	coll := s.db.Collection("galleries")

	cursor, err := coll.Find(ctx, bson.M{"collectionId": collectionId, "userId": userId})
	if err != nil {
		return nil, err
	}

	result := make([]domain.GalleryDB, 0)
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MongoGallery) GetGallery(ctx context.Context, galleryId primitive.ObjectID, userId string) (domain.GalleryDB, error) {
	coll := s.db.Collection("galleries")

	var result domain.GalleryDB
	err := coll.FindOne(ctx, bson.M{"_id": galleryId, "userId": userId}).Decode(&result)
	if err != nil {
		return domain.GalleryDB{}, err
	}

	return result, nil
}

func (s *MongoGallery) CreateGallery(ctx context.Context, collectionId primitive.ObjectID, name, userId string) (string, error) {

	galleriesColl := s.db.Collection("galleries")
	galleryID := primitive.NewObjectID()

	gallery := bson.D{
		{"_id", galleryID},
		{"collectionId", collectionId},
		{"name", name},
		{"userId", userId},
		{"createdAt", primitive.NewDateTimeFromTime(time.Now().UTC())},
		{"updatedAt", primitive.NewDateTimeFromTime(time.Now().UTC())},
		{"sharing", bson.D{
			{"sharingEnabled", false},
		}},
	}
	_, err := galleriesColl.InsertOne(ctx, gallery)
	if err != nil {
		return "", err
	}

	return galleryID.Hex(), nil
}

func (s *MongoGallery) DeleteGallery(ctx context.Context, galleryId primitive.ObjectID, userId string) error {
	coll := s.db.Collection("galleries")
	_, err := coll.DeleteOne(ctx, bson.M{"_id": galleryId, "userId": userId})
	return err
}

func (s *MongoGallery) UpdateGallery(ctx context.Context, galleryId primitive.ObjectID, userId string, opts ...domain.GalleryUpdateOption) (domain.GalleryDB, error) {

	updateOptions := &domain.GalleryUpdateOptions{
		SetFields: bson.D{},
	}
	for _, opt := range opts {
		opt(updateOptions)
	}

	coll := s.db.Collection("galleries")
	filter := bson.M{"_id": galleryId, "userId": userId}
	update := bson.D{
		{"$set", updateOptions.SetFields},
		{"$currentDate", bson.D{
			{"updatedAt", true},
		}},
	}
	findOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var gallery domain.GalleryDB
	err := coll.FindOneAndUpdate(ctx, filter, update, findOpts).Decode(&gallery)
	return gallery, err
}
