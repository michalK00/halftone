package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type GalleryDB struct {
	ID primitive.ObjectID `bson:"_id" json:"id"` 
	Name string `bson:"name" json:"name"`
	ExpiryDate primitive.DateTime `bson:"expiry_date" json:"expiry_date"`
}