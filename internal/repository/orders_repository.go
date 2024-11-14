package repository

import "go.mongodb.org/mongo-driver/mongo"

type MongoOrder struct {
	db *mongo.Database
}

func NewMongoOrder(db *mongo.Database) *MongoGallery {
	return &MongoGallery{
		db: db,
	}
}
