package repository

import (
	"context"
	"github.com/michalK00/sg-qr/internal/domain"
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

func (s *MongoPhoto) PhotoExists(ctx context.Context, photoId primitive.ObjectID) (bool, error) {
	coll := s.db.Collection("photos")

	count, err := coll.CountDocuments(ctx, bson.D{{"_id", photoId}}, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *MongoPhoto) GalleryPhotoCount(ctx context.Context, galleryId primitive.ObjectID) (int64, error) {
	coll := s.db.Collection("photos")
	count, err := coll.CountDocuments(ctx, bson.D{{"galleryId", galleryId}})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *MongoPhoto) GetPhotos(ctx context.Context, galleryId primitive.ObjectID) ([]domain.PhotoDB, error) {
	coll := s.db.Collection("photos")

	var result []domain.PhotoDB
	// returns only uploaded
	cursor, err := coll.Find(ctx, bson.D{{"galleryId", galleryId}, {"status", 1}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MongoPhoto) GetPhoto(ctx context.Context, galleryId primitive.ObjectID) (domain.PhotoDB, error) {
	coll := s.db.Collection("photos")

	var result domain.PhotoDB
	err := coll.FindOne(ctx, bson.D{{"_id", galleryId}}).Decode(&result)
	if err != nil {
		return domain.PhotoDB{}, err
	}

	return result, nil
}

func (s *MongoPhoto) CreatePhoto(ctx context.Context, collectionId primitive.ObjectID, galleryId primitive.ObjectID, originalFilename string) (primitive.ObjectID, error) {

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
		{"originalFilename", originalFilename},
		{"createdAt", primitive.NewDateTimeFromTime(time.Now())},
		{"updatedAt", primitive.NewDateTimeFromTime(time.Now())},
		{"status", "pending"},
		{"objectKey", path.Join(collectionId.Hex(), galleryId.Hex(), "photos", photoId.Hex()+ext)},
	}
	_, err := coll.InsertOne(ctx, photo)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return photoId, nil
}

func (s *MongoPhoto) CreatePhotos(ctx context.Context, collectionId primitive.ObjectID, galleryId primitive.ObjectID, originalFilenames []string) ([]primitive.ObjectID, error) {

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
			{"originalFilename", filename},
			{"createdAt", primitive.NewDateTimeFromTime(time.Now())},
			{"updatedAt", primitive.NewDateTimeFromTime(time.Now())},
			{"status", "pending"},
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

func (s *MongoPhoto) DeletePhoto(ctx context.Context, photoId primitive.ObjectID) error {
	coll := s.db.Collection("photos")
	_, err := coll.DeleteOne(ctx, bson.D{{"_id", photoId}})
	return err
}

func (s *MongoPhoto) DeletePhotos(ctx context.Context, photoIds []primitive.ObjectID) error {
	coll := s.db.Collection("photos")
	filter := bson.M{"_id": bson.M{"$in": photoIds}}
	_, err := coll.DeleteMany(ctx, filter)
	return err
}

func (s *MongoPhoto) UpdatePhoto(ctx context.Context, photoId primitive.ObjectID, status domain.PhotoStatus) (domain.PhotoDB, error) {
	coll := s.db.Collection("photos")
	filter := bson.D{{"_id", photoId}}
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
