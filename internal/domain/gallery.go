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
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	CollectionId      primitive.ObjectID `bson:"collectionId" json:"collectionId"`
	Name              string             `bson:"name" json:"name"`
	CreatedAt         primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt         primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	AccessToken       string             `bson:"accessToken" json:"accessToken"`
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

func WithAccessToken(accessToken string) GalleryUpdateOption {
	return func(opts *GalleryUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "accessToken", Value: accessToken})
	}
}

func WithSharingEnabled(enabled bool) GalleryUpdateOption {
	return func(opts *GalleryUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "enabled", Value: enabled})
	}
}

func WithSharingExpiryDate(date primitive.DateTime) GalleryUpdateOption {
	return func(opts *GalleryUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "sharingExpiryDate", Value: date})
	}
}

func WithValidatedSharingExpiryDate(date primitive.DateTime) GalleryUpdateOption {
	return func(opts *GalleryUpdateOptions) {
		if date < primitive.NewDateTimeFromTime(time.Now()) {
			date = primitive.NewDateTimeFromTime(time.Now().Add(time.Hour * 24 * 7))
		}
		opts.SetFields = append(opts.SetFields, bson.E{Key: "sharingExpiryDate", Value: date})
	}
}
