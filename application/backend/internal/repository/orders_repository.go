package repository

import (
	"context"
	"errors"
	"github.com/michalK00/halftone/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoOrder struct {
	db *mongo.Database
}

func NewMongoOrder(db *mongo.Database) *MongoOrder {
	return &MongoOrder{
		db: db,
	}
}

func (s *MongoOrder) OrderExists(ctx context.Context, orderId primitive.ObjectID, userId string) (bool, error) {
	// Check if order exists and belongs to a gallery owned by the user
	pipeline := []bson.M{
		{"$match": bson.M{"_id": orderId}},
		{"$lookup": bson.M{
			"from":         "galleries",
			"localField":   "gallery_id",
			"foreignField": "_id",
			"as":           "gallery",
		}},
		{"$unwind": "$gallery"},
		{"$match": bson.M{"gallery.userId": userId}},
		{"$limit": 1},
	}

	coll := s.db.Collection("orders")
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	return cursor.Next(ctx), nil
}

func (s *MongoOrder) GetOrders(ctx context.Context, userId string) ([]domain.OrderDB, error) {
	// Get all orders for galleries owned by the user
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "galleries",
			"localField":   "gallery_id",
			"foreignField": "_id",
			"as":           "gallery",
		}},
		{"$unwind": "$gallery"},
		{"$match": bson.M{"gallery.userId": userId}},
		{"$project": bson.M{
			"gallery": 0, // Remove gallery data from result
		}},
		{"$sort": bson.M{"created_at": -1}}, // Most recent first
	}

	coll := s.db.Collection("orders")
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.OrderDB
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *MongoOrder) GetOrder(ctx context.Context, orderId primitive.ObjectID, userId string) (domain.OrderDB, error) {
	// Get order only if it belongs to a gallery owned by the user
	pipeline := []bson.M{
		{"$match": bson.M{"_id": orderId}},
		{"$lookup": bson.M{
			"from":         "galleries",
			"localField":   "gallery_id",
			"foreignField": "_id",
			"as":           "gallery",
		}},
		{"$unwind": "$gallery"},
		{"$match": bson.M{"gallery.userId": userId}},
		{"$project": bson.M{
			"gallery": 0, // Remove gallery data from result
		}},
	}

	coll := s.db.Collection("orders")
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return domain.OrderDB{}, err
	}
	defer cursor.Close(ctx)

	var order domain.OrderDB
	if cursor.Next(ctx) {
		if err := cursor.Decode(&order); err != nil {
			return domain.OrderDB{}, err
		}
		return order, nil
	}

	return domain.OrderDB{}, mongo.ErrNoDocuments
}

func (s *MongoOrder) CreateOrder(ctx context.Context, galleryId primitive.ObjectID, clientEmail, comment string, photoIds []primitive.ObjectID) (string, error) {
	ordersColl := s.db.Collection("orders")
	orderID := primitive.NewObjectID()

	// Convert photo IDs to OrderPhoto structs
	photos := make([]domain.OrderPhoto, len(photoIds))
	for i, photoId := range photoIds {
		photos[i] = domain.OrderPhoto{PhotoID: photoId}
	}

	order := domain.OrderDB{
		ID:          orderID,
		GalleryID:   galleryId,
		ClientEmail: clientEmail,
		Comment:     comment,
		Status:      domain.OrderStatusPending,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Photos:      photos,
	}

	_, err := ordersColl.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}

	return orderID.Hex(), nil
}

func (s *MongoOrder) UpdateOrder(ctx context.Context, orderId primitive.ObjectID, userId string, opts ...domain.OrderUpdateOption) (domain.OrderDB, error) {
	updateOptions := &domain.OrderUpdateOptions{
		SetFields: bson.D{},
	}
	for _, opt := range opts {
		opt(updateOptions)
	}

	// First verify the order belongs to user's gallery
	exists, err := s.OrderExists(ctx, orderId, userId)
	if err != nil {
		return domain.OrderDB{}, err
	}
	if !exists {
		return domain.OrderDB{}, mongo.ErrNoDocuments
	}

	coll := s.db.Collection("orders")
	filter := bson.M{"_id": orderId}
	update := bson.D{
		{"$set", updateOptions.SetFields},
		{"$currentDate", bson.D{
			{"updated_at", true},
		}},
	}

	findOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var order domain.OrderDB
	err = coll.FindOneAndUpdate(ctx, filter, update, findOpts).Decode(&order)
	return order, err
}

func (s *MongoOrder) DeleteOrder(ctx context.Context, orderId primitive.ObjectID, userId string) error {
	// First verify the order belongs to user's gallery
	exists, err := s.OrderExists(ctx, orderId, userId)
	if err != nil {
		return err
	}
	if !exists {
		return mongo.ErrNoDocuments
	}

	coll := s.db.Collection("orders")
	_, err = coll.DeleteOne(ctx, bson.M{"_id": orderId})
	return err
}

func (s *MongoOrder) OrderExistsForGallery(ctx context.Context, galleryId primitive.ObjectID) (bool, error) {
	coll := s.db.Collection("orders")
	filter := bson.M{"gallery_id": galleryId}
	var order domain.OrderDB
	if err := coll.FindOne(ctx, filter).Decode(&order); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}
