package repository

import (
	"context"
	"github.com/michalK00/halftone/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"path"
	"path/filepath"
	"time"
)

type MongoPhoto struct {
	db *mongo.Database
}

func NewMongoPhoto(db *mongo.Database) *MongoPhoto {
	return &MongoPhoto{
		db: db,
	}
}

func (s *MongoPhoto) PhotoExists(ctx context.Context, photoId primitive.ObjectID, userId string) (bool, error) {
	coll := s.db.Collection("photos")

	count, err := coll.CountDocuments(ctx, bson.M{"_id": photoId, "userId": userId}, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *MongoPhoto) GalleryPhotoCount(ctx context.Context, galleryId primitive.ObjectID, userId string) (int64, error) {
	coll := s.db.Collection("photos")
	count, err := coll.CountDocuments(ctx, bson.M{"galleryId": galleryId, "userId": userId})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *MongoPhoto) GetPhotos(ctx context.Context, galleryId primitive.ObjectID, userId string) ([]domain.PhotoDB, error) {
	coll := s.db.Collection("photos")

	var result []domain.PhotoDB
	// returns only uploaded and shared
	cursor, err := coll.Find(ctx, bson.M{
		"galleryId": galleryId,
		"status":    bson.D{{"$in", primitive.A{1, 2}}},
		"userId":    userId,
	})

	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MongoPhoto) GetPhoto(ctx context.Context, photoId primitive.ObjectID, userId string) (domain.PhotoDB, error) {
	coll := s.db.Collection("photos")

	var result domain.PhotoDB
	err := coll.FindOne(ctx, bson.M{"_id": photoId, "userId": userId}).Decode(&result)
	if err != nil {
		return domain.PhotoDB{}, err
	}

	return result, nil
}

func (s *MongoPhoto) CreatePhoto(ctx context.Context, collectionId primitive.ObjectID, galleryId primitive.ObjectID, originalFilename string, userId string) (primitive.ObjectID, error) {

	coll := s.db.Collection("photos")
	photoId := primitive.NewObjectID()
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".jpg"
	}

	photo := bson.D{
		{"_id", photoId},
		{"collectionId", collectionId},
		{"galleryId", galleryId},
		{"userId", userId},
		{"originalFilename", originalFilename},
		{"createdAt", primitive.NewDateTimeFromTime(time.Now().UTC())},
		{"updatedAt", primitive.NewDateTimeFromTime(time.Now().UTC())},
		{"status", "pending"},
		{"objectKey", path.Join(collectionId.Hex(), galleryId.Hex(), "photos", photoId.Hex()+ext)},
	}
	_, err := coll.InsertOne(ctx, photo)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return photoId, nil
}

func (s *MongoPhoto) CreatePhotos(ctx context.Context, collectionId primitive.ObjectID, galleryId primitive.ObjectID, originalFilenames []string, userId string) ([]primitive.ObjectID, error) {

	coll := s.db.Collection("photos")

	photoIds := make([]primitive.ObjectID, len(originalFilenames))
	documents := make([]interface{}, len(originalFilenames))
	for i, filename := range originalFilenames {
		photoId := primitive.NewObjectID()
		ext := filepath.Ext(filename)
		if ext == "" {
			ext = ".jpg"
		}
		photo := bson.D{
			{"_id", photoId},
			{"collectionId", collectionId},
			{"galleryId", galleryId},
			{"userId", userId},
			{"originalFilename", filename},
			{"createdAt", primitive.NewDateTimeFromTime(time.Now().UTC())},
			{"updatedAt", primitive.NewDateTimeFromTime(time.Now().UTC())},
			{"status", domain.PhotoStatus(0)},
			{"objectKey", path.Join(collectionId.Hex(), galleryId.Hex(), "photos", photoId.Hex()+ext)},
		}
		documents[i] = photo
		photoIds[i] = photoId
	}

	_, err := coll.InsertMany(ctx, documents)
	if err != nil {
		return []primitive.ObjectID(nil), err
	}

	return photoIds, nil
}

func (s *MongoPhoto) SoftDeletePhoto(ctx context.Context, photoId primitive.ObjectID, userId string) error {
	coll := s.db.Collection("photos")
	filter := bson.M{"_id": photoId, "userId": userId}
	update := bson.D{
		{"$set", bson.D{
			{"status", domain.PhotoStatus(3)},
		}},
		{"$currentDate", bson.D{
			{"updatedAt", true},
		}},
	}
	opts := options.FindOneAndUpdate()
	return coll.FindOneAndUpdate(ctx, filter, update, opts).Err()
}

func (s *MongoPhoto) DeletePhoto(ctx context.Context, photoId primitive.ObjectID, userId string) error {
	coll := s.db.Collection("photos")
	_, err := coll.DeleteOne(ctx, bson.M{"_id": photoId, "userId": userId})
	return err
}

func (s *MongoPhoto) DeletePhotos(ctx context.Context, photoIds []primitive.ObjectID, userId string) error {
	coll := s.db.Collection("photos")
	filter := bson.M{"_id": bson.M{"$in": photoIds}, "userId": userId}
	_, err := coll.DeleteMany(ctx, filter)
	return err
}

func (s *MongoPhoto) UpdatePhoto(ctx context.Context, photoId primitive.ObjectID, status domain.PhotoStatus, userId string) (domain.PhotoDB, error) {
	coll := s.db.Collection("photos")
	filter := bson.M{"_id": photoId, "userId": userId}
	update := bson.D{
		{"$set", bson.D{
			{"status", status},
		}},
		{"$currentDate", bson.D{
			{"updatedAt", true},
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var photo domain.PhotoDB
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&photo)
	return photo, err
}
