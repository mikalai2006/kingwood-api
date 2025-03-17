package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArchiveMessage struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`
	// ProductID primitive.ObjectID     `json:"productId" bson:"productId"`
	OrderID        primitive.ObjectID     `json:"orderId" bson:"orderId"`
	Status         int                    `json:"status" bson:"status"`
	ArchiveMessage string                 `json:"ArchiveMessage" bson:"ArchiveMessage"`
	Props          map[string]interface{} `json:"props" bson:"props"`

	Images []MessageImage `json:"images" bson:"images"`
	// Statuses []ArchiveMessageStatus `json:"statuses" bson:"statuses"`
	// User User `json:"user,omitempty" bson:"user,omitempty"`
	// Images []ArchiveMessageImage `json:"images" bson:"images,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`

	Meta ArchiveMeta `json:"meta" bson:"meta"`
}

type ArchiveMessageInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`
	// ProductID primitive.ObjectID     `json:"productId" bson:"productId"`
	OrderID primitive.ObjectID `json:"orderId" bson:"orderId"`
	// Status  int                    `json:"status" bson:"status"`
	ArchiveMessage string                 `json:"ArchiveMessage" bson:"ArchiveMessage"`
	Props          map[string]interface{} `json:"props" bson:"props"`

	Images []MessageImage `json:"images" bson:"images"`
	// Statuses []ArchiveMessageStatus `json:"statuses" bson:"statuses"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`

	Meta ArchiveMeta `json:"meta" bson:"meta"`
}

type ArchiveMessageFilter struct {
	ID     string `json:"id,omitempty"`
	UserID string `json:"userId,omitempty"`
	// ProductID *primitive.ObjectID        `json:"productId,omitempty"`
	OrderID []string            `json:"orderId" bson:"orderId"`
	Sort    []*FilterSortParams `json:"$sort,omitempty"`
	Limit   *int                `json:"$limit,omitempty"`
	Skip    *int                `json:"$skip,omitempty"`
}
