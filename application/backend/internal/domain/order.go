package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
)

type OrderDB struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	GalleryID   primitive.ObjectID `bson:"gallery_id" json:"galleryId"`
	ClientEmail string             `bson:"client_email" json:"clientEmail"`
	Comment     string             `bson:"comment" json:"comment"`
	Status      OrderStatus        `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	Photos      []OrderPhoto       `bson:"photos" json:"photos"`
}

type OrderPhoto struct {
	PhotoID primitive.ObjectID `bson:"photo_id" json:"photoId"`
}

type OrderRepository interface {
	// User endpoints
	GetOrders(ctx context.Context, userId string) ([]OrderDB, error)
	GetOrder(ctx context.Context, orderId primitive.ObjectID, userId string) (OrderDB, error)
	UpdateOrder(ctx context.Context, orderId primitive.ObjectID, userId string, opts ...OrderUpdateOption) (OrderDB, error)
	DeleteOrder(ctx context.Context, orderId primitive.ObjectID, userId string) error

	// Client endpoints
	CreateOrder(ctx context.Context, galleryId primitive.ObjectID, clientEmail, comment string, photoIds []primitive.ObjectID) (string, error)

	// Helper methods
	OrderExists(ctx context.Context, orderId primitive.ObjectID, userId string) (bool, error)
	OrderExistsForGallery(ctx context.Context, galleryId primitive.ObjectID) (bool, error)
}

type OrderUpdateOption func(*OrderUpdateOptions)

type OrderUpdateOptions struct {
	SetFields bson.D
}

func WithOrderStatus(status OrderStatus) OrderUpdateOption {
	return func(opts *OrderUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "status", Value: status})
	}
}

func WithOrderComment(comment string) OrderUpdateOption {
	return func(opts *OrderUpdateOptions) {
		opts.SetFields = append(opts.SetFields, bson.E{Key: "comment", Value: comment})
	}
}
