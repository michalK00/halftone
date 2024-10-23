package galleries

import (
	"context"
	"github.com/michalK00/sg-qr/app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type GalleryStorage struct {
	db *mongo.Database
}

func NewGalleryStorage(db *mongo.Database) *GalleryStorage {
	return &GalleryStorage{
		db: db,
	}
}

func (s *GalleryStorage) galleryExists(ctx context.Context, galleryId primitive.ObjectID) (bool, error) {
	coll := s.db.Collection("galleries")

	count, err := coll.CountDocuments(ctx, bson.D{{"_id", galleryId}}, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *GalleryStorage) getGalleries(ctx context.Context, collectionId primitive.ObjectID) ([]domain.GalleryDB, error) {
	coll := s.db.Collection("galleries")

	var result []domain.GalleryDB
	cursor, err := coll.Find(ctx, bson.D{{"collectionId", collectionId}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *GalleryStorage) getGallery(ctx context.Context, galleryId primitive.ObjectID) (domain.GalleryDB, error) {
	coll := s.db.Collection("galleries")

	var result domain.GalleryDB
	err := coll.FindOne(ctx, bson.D{{"_id", galleryId}}).Decode(&result)
	if err != nil {
		return domain.GalleryDB{}, err
	}

	return result, nil
}

func (s *GalleryStorage) createGallery(ctx context.Context, collectionId primitive.ObjectID, name string) (string, error) {

	mongoCollection := s.db.Collection("galleries")
	galleryID := primitive.NewObjectID()

	gallery := bson.D{
		{"_id", galleryID},
		{"collectionId", collectionId},
		{"name", name},
		{"createdAt", primitive.NewDateTimeFromTime(time.Now())},
		{"updatedAt", primitive.NewDateTimeFromTime(time.Now())},
		{"sharingEnabled", false},
	}
	_, err := mongoCollection.InsertOne(ctx, gallery)
	if err != nil {
		return "", err
	}

	return galleryID.Hex(), nil
}

func (s *GalleryStorage) deleteGallery(ctx context.Context, galleryId primitive.ObjectID) error {
	coll := s.db.Collection("galleries")
	_, err := coll.DeleteOne(ctx, bson.D{{"_id", galleryId}})
	return err
}

func (s *GalleryStorage) updateGallery(ctx context.Context, galleryId primitive.ObjectID, name string, sharingEnabled bool, sharingExpiryDate primitive.DateTime) error {
	coll := s.db.Collection("galleries")
	filter := bson.D{{"_id", galleryId}}
	update := bson.D{
		{"$set", bson.D{
			{"name", name},
			{"sharingEnabled", sharingEnabled},
			{"sharingExpiryDate", sharingExpiryDate},
		}},
		{"$currentDate", bson.D{
			{"updatedAt", true},
		}},
	}
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}
