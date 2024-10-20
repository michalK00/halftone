package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type CollectionDB struct {
	ID primitive.ObjectID `bson:"_id" json:"id"` 
	Name string `bson:"name" json:"name"`
	Galleries []GalleryDB `bson:"galleries" json:"galleries"`
}