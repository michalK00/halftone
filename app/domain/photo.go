package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type PhotoDB struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	GalleryId    primitive.ObjectID `bson:"galleryId" json:"galleryId"`
	CollectionId primitive.ObjectID `bson:"collectionId" json:"collectionId"`
	Status       string             `bson:"status" json:"status"`
	CreatedAt    primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt    primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	Metadata     PhotoMetadata      `bson:"metadata" json:"metadata"`
}

type PhotoMetadata struct {
	Size             int64  `bson:"size" json:"size"`
	OriginalFilename string `bson:"originalFilename" json:"originalFilename"`
}
