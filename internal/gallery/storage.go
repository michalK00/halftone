package gallery

import (
	"context"
	"time"

	"github.com/michalK00/sg-qr/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GalleryStorage struct {
	db *mongo.Database
}

func NewGalleryStorage(db *mongo.Database) *GalleryStorage {
	return &GalleryStorage{
		db: db,
	}
}

func (s* GalleryStorage) getAllGalleries(ctx context.Context, collectionId primitive.ObjectID) ([]domain.GalleryDB, error) {
	mongoCollection := s.db.Collection("collections")

	var result domain.CollectionDB
	err := mongoCollection.FindOne(ctx, bson.M{"_id": collectionId}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Galleries, nil
}

func (s* GalleryStorage) createGallery(ctx context.Context, collectionId primitive.ObjectID, name string, expiryDate time.Time) (string, error) {
	
	mongoCollection := s.db.Collection("collections")
	galleryID := primitive.NewObjectID()

	gallery := bson.M{
		"_id": galleryID,
		"name": name,
		"expiry_date": primitive.NewDateTimeFromTime(expiryDate),
	}

	result, err := mongoCollection.UpdateOne(ctx, bson.M{"_id": collectionId}, bson.M{"$push": bson.M{"galleries": gallery}})

	if err != nil {
		return "", err
	}
	if result.MatchedCount == 0 {
		return "", err
	}
	
	return galleryID.Hex(), nil;
}