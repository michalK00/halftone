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
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	CollectionId primitive.ObjectID `bson:"collectionId" json:"collectionId"`
	Name         string             `bson:"name" json:"name"`
	UserId       string             `bson:"userId" json:"userId"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	Sharing      Sharing            `bson:"sharing" json:"sharing"`
	PhotoOptions PhotoOptions       `bson:"photoOptions" json:"photoOptions"`
}

type Sharing struct {
	SharingEnabled    bool      `bson:"sharingEnabled" json:"sharingEnabled"`
	SharingExpiryDate time.Time `bson:"sharingExpiryDate" json:"sharingExpiryDate"`
	AccessToken       string    `bson:"accessToken" json:"accessToken"`
	SharingUrl        string    `bson:"sharingUrl" json:"sharingUrl"`
}

type PhotoOptions struct {
	Downsize  bool `bson:"downsize" json:"downsize"`
	Watermark bool `bson:"watermark" json:"watermark"`
}

type GalleryRepository interface {
	GetGalleryByID(ctx context.Context, galleryId primitive.ObjectID) (GalleryDB, error)
	GalleryExists(ctx context.Context, galleryId primitive.ObjectID, userId string) (bool, error)
	CollectionGalleryCount(ctx context.Context, collectionId primitive.ObjectID, userId string) (int64, error)
	GetGalleries(ctx context.Context, collectionId primitive.ObjectID, userId string) ([]GalleryDB, error)
	GetGallery(ctx context.Context, galleryId primitive.ObjectID, userId string) (GalleryDB, error)
	CreateGallery(ctx context.Context, collectionId primitive.ObjectID, name, userId string) (string, error)
	DeleteGallery(ctx context.Context, galleryId primitive.ObjectID, userId string) error
	UpdateGallery(ctx context.Context, galleryId primitive.ObjectID, userId string, opts ...GalleryUpdateOption) (GalleryDB, error)
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
			{Key: "sharingEnabled", Value: sharing.SharingEnabled},
			{Key: "accessToken", Value: sharing.AccessToken},
			{Key: "sharingExpiryDate", Value: sharing.SharingExpiryDate},
			{Key: "sharingUrl", Value: sharing.SharingUrl},
		}})
	}
}
