package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

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
