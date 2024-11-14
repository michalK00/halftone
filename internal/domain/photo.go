package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PhotoDB struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	GalleryId        primitive.ObjectID `bson:"galleryId" json:"galleryId"`
	CollectionId     primitive.ObjectID `bson:"collectionId" json:"collectionId"`
	Status           PhotoStatus        `bson:"status" json:"status"`
	CreatedAt        primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt        primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	OriginalFilename string             `bson:"originalFilename" json:"originalFilename"`
	ObjectKey        string             `bson:"objectKey" json:"objectKey"`
}

type PhotoStatus int64

const (
	Pending PhotoStatus = iota
	Uploaded
)

type PhotoRepository interface {
	PhotoExists(ctx context.Context, photoId primitive.ObjectID) (bool, error)
	GalleryPhotoCount(ctx context.Context, galleryId primitive.ObjectID) (int64, error)
	GetPhotos(ctx context.Context, galleryId primitive.ObjectID) ([]PhotoDB, error)
	GetPhoto(ctx context.Context, galleryId primitive.ObjectID) (PhotoDB, error)
	CreatePhoto(ctx context.Context, collectionId primitive.ObjectID, galleryId primitive.ObjectID, originalFilename string) (primitive.ObjectID, error)
	CreatePhotos(ctx context.Context, collectionId primitive.ObjectID, galleryId primitive.ObjectID, originalFilename []string) ([]primitive.ObjectID, error)
	DeletePhoto(ctx context.Context, photoId primitive.ObjectID) error
	DeletePhotos(ctx context.Context, photoIds []primitive.ObjectID) error
	UpdatePhoto(ctx context.Context, photoId primitive.ObjectID, status PhotoStatus) (PhotoDB, error)
}
