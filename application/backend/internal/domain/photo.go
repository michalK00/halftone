package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PhotoDB struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	GalleryId          primitive.ObjectID `bson:"galleryId" json:"galleryId"`
	CollectionId       primitive.ObjectID `bson:"collectionId" json:"collectionId"`
	UserId             string             `bson:"userId" json:"userId"`
	Status             PhotoStatus        `bson:"status" json:"status"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
	OriginalFilename   string             `bson:"originalFilename" json:"originalFilename"`
	ObjectKey          string             `bson:"objectKey" json:"objectKey"`
	ClientObjectKey    string             `bson:"clientObjectKey" json:"clientObjectKey"`
	ThumbnailObjectKey string             `bson:"thumbnailObjectKey" json:"thumbnailObjectKey"`
}

type PhotoStatus int64

const (
	Pending PhotoStatus = iota
	Uploaded
	Shared
	PendingDeletion
)

type PhotoRepository interface {
	PhotoExists(ctx context.Context, photoId primitive.ObjectID, userId string) (bool, error)
	GalleryPhotoCount(ctx context.Context, galleryId primitive.ObjectID, userId string) (int64, error)
	GetPhotos(ctx context.Context, galleryId primitive.ObjectID, userId string) ([]PhotoDB, error)
	GetPhoto(ctx context.Context, photoId primitive.ObjectID, userId string) (PhotoDB, error)
	CreatePhoto(ctx context.Context, collectionId primitive.ObjectID, galleryId primitive.ObjectID, originalFilename string, userId string) (primitive.ObjectID, error)
	CreatePhotos(ctx context.Context, collectionId primitive.ObjectID, galleryId primitive.ObjectID, originalFilename []string, userId string) ([]primitive.ObjectID, error)
	DeletePhoto(ctx context.Context, photoId primitive.ObjectID, userId string) error
	SoftDeletePhoto(ctx context.Context, photoId primitive.ObjectID, userId string) error
	DeletePhotos(ctx context.Context, photoIds []primitive.ObjectID, userId string) error
	UpdatePhoto(ctx context.Context, photoId primitive.ObjectID, status PhotoStatus, userId string) (PhotoDB, error)
	GetSharedPhotosByGallery(ctx context.Context, galleryId primitive.ObjectID) ([]PhotoDB, error)
	GetSharedPhotoById(ctx context.Context, photoId primitive.ObjectID) (PhotoDB, error)
	VerifyPhotosInGallery(ctx context.Context, galleryId primitive.ObjectID, photoIds []primitive.ObjectID) (bool, error)
}
