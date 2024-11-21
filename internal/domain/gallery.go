package domain

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GalleryDB struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	CollectionId   primitive.ObjectID `bson:"collectionId" json:"collectionId"`
	Name           string             `bson:"name" json:"name"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	SharingOptions SharingOptions     `bson:"sharingOptions" json:"sharingOptions"`
}

type SharingOptions struct {
	SharingEnabled    bool               `bson:"sharingEnabled" json:"sharingEnabled"`
	AccessToken       string             `bson:"accessToken" json:"accessToken"`
	SharingExpiryDate time.Time          `bson:"sharingExpiryDate" json:"sharingExpiryDate"`
	SharingUrl        string             `bson:"sharingUrl" json:"sharingUrl"`
	SharingCleanupJob primitive.ObjectID `bson:"sharingCleanupJob" json:"sharingCleanupJob"`
}

type GalleryRepository interface {
	GalleryExists(ctx context.Context, galleryId primitive.ObjectID) (bool, error)
	CollectionGalleryCount(ctx context.Context, collectionId primitive.ObjectID) (int64, error)
	GetGalleries(ctx context.Context, collectionId primitive.ObjectID) ([]GalleryDB, error)
	GetGallery(ctx context.Context, galleryId primitive.ObjectID) (GalleryDB, error)
	CreateGallery(ctx context.Context, collectionId primitive.ObjectID, name string) (string, error)
	DeleteGallery(ctx context.Context, galleryId primitive.ObjectID) error
	UpdateGallery(ctx context.Context, galleryId primitive.ObjectID, opts ...GalleryUpdateOption) (GalleryDB, error)
}

func GenerateAccessToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

type GalleryUpdateOption func(*GalleryUpdateOptions)

type GalleryUpdateOptions struct {
	SetFields bson.D
}

func WithName(name string) GalleryUpdateOption {
	return func(opts *GalleryUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "name", Value: name})
	}
}

func WithSharingOptions(sharingOptions SharingOptions) GalleryUpdateOption {
	return func(opts *GalleryUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "sharingOptions", Value: bson.D{
			{Key: "sharingEnabled", Value: sharingOptions.SharingEnabled},
			{Key: "accessToken", Value: sharingOptions.AccessToken},
			{Key: "sharingExpiryDate", Value: sharingOptions.SharingExpiryDate},
			{Key: "sharingUrl", Value: sharingOptions.SharingUrl},
			{Key: "sharingCleanupJob", Value: sharingOptions.SharingCleanupJob},
		}})
	}
}
func WithSharingEnabled(sharingEnabled bool) GalleryUpdateOption {
	return func(opts *GalleryUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "sharingOptions", Value: bson.D{
			{Key: "sharingEnabled", Value: sharingEnabled},
		}})
	}
}
