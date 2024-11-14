package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GalleryDB struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	CollectionId      primitive.ObjectID `bson:"collectionId" json:"collectionId"`
	Name              string             `bson:"name" json:"name"`
	CreatedAt         primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt         primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	PasswordHash      string             `bson:"passwordHash" json:"passwordHash"`
	PasswordAttempts  int                `bson:"passwordAttempts" json:"passwordAttempts"`
	SharingEnabled    bool               `bson:"sharingEnabled" json:"sharingEnabled"`
	SharingExpiryDate primitive.DateTime `bson:"sharingExpiryDate" json:"sharingExpiryDate"`
}

type SharingOptions struct {
	Watermark bool `bson:"watermark" json:"watermark"`
	Downsize  bool `bson:"downsize" json:"downsize"`
}

type GalleryRepository interface {
	GalleryExists(ctx context.Context, galleryId primitive.ObjectID) (bool, error)
	CollectionGalleryCount(ctx context.Context, collectionId primitive.ObjectID) (int64, error)
	GetGalleries(ctx context.Context, collectionId primitive.ObjectID) ([]GalleryDB, error)
	GetGallery(ctx context.Context, galleryId primitive.ObjectID) (GalleryDB, error)
	CreateGallery(ctx context.Context, collectionId primitive.ObjectID, name string) (string, error)
	DeleteGallery(ctx context.Context, galleryId primitive.ObjectID) error
	UpdateGallery(ctx context.Context, galleryId primitive.ObjectID, name string, sharingEnabled bool, sharingExpiryDate primitive.DateTime) (GalleryDB, error)
}
