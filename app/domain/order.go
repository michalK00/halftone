package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderDB struct {
	ID           primitive.ObjectID   `bson:"_id" json:"id"`
	GalleryId    primitive.ObjectID   `bson:"galleryId" json:"galleryId"`
	CollectionId primitive.ObjectID   `bson:"collectionId" json:"collectionId"`
	Status       string               `bson:"status" json:"status"`
	CreatedAt    primitive.DateTime   `bson:"createdAt" json:"createdAt"`
	UpdatedAt    primitive.DateTime   `bson:"updatedAt" json:"updatedAt"`
	Photos       []primitive.ObjectID `bson:"photos" json:"photos"`
}
