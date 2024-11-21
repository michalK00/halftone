package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderDB struct {
	ID           primitive.ObjectID   `bson:"_id" json:"id"`
	GalleryId    primitive.ObjectID   `bson:"galleryId" json:"galleryId"`
	CollectionId primitive.ObjectID   `bson:"collectionId" json:"collectionId"`
	Status       string               `bson:"status" json:"status"`
	CreatedAt    time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time            `bson:"updatedAt" json:"updatedAt"`
	Photos       []primitive.ObjectID `bson:"photos" json:"photos"`
}

type OrderRepository interface {
}
