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
	Sharing        Sharing            `bson:"sharing" json:"sharing"`
	SharingOptions SharingOptions     `bson:"sharingOptions" json:"sharingOptions"`
}

type Sharing struct {
	SharingExpiryDate time.Time `bson:"sharingExpiryDate" json:"sharingExpiryDate"`
	AccessToken       string    `bson:"accessToken" json:"accessToken"`
	SharingUrl        string    `bson:"sharingUrl" json:"sharingUrl"`
}

type SharingOptions struct {
	Downsize  bool `bson:"downsize" json:"downsize"`
	Watermark bool `bson:"watermark" json:"watermark"`
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

func WithSharing(sharing Sharing) GalleryUpdateOption {
	return func(opts *GalleryUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "sharing", Value: bson.D{
			{Key: "accessToken", Value: sharing.AccessToken},
			{Key: "sharingExpiryDate", Value: sharing.SharingExpiryDate},
			{Key: "sharingUrl", Value: sharing.SharingUrl},
		}})
	}
}
