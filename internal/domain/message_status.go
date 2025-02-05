package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageStatus struct {
	ID      primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID  primitive.ObjectID     `json:"userId" bson:"userId"`
	OrderID primitive.ObjectID     `json:"orderId" bson:"orderId"`
	Status  *int                   `json:"status" bson:"status"`
	Props   map[string]interface{} `json:"props" bson:"props"`

	User User `json:"user,omitempty" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type MessageStatusMongo struct {
	ID      primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID  primitive.ObjectID     `json:"userId" bson:"userId"`
	OrderID primitive.ObjectID     `json:"orderId" bson:"orderId"`
	Status  *int                   `json:"status" bson:"status"`
	Props   map[string]interface{} `json:"props" bson:"props"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type MessageStatusFilter struct {
	ID      *primitive.ObjectID `json:"id" bson:"id"`
	UserID  *primitive.ObjectID `json:"userId" bson:"userId"`
	OrderID *primitive.ObjectID `json:"orderId" bson:"orderId"`
	Status  *int                `json:"status" bson:"status"`
	Sort    []*FilterSortParams `json:"sort" bson:"sort"`
	Limit   *int                `json:"limit" bson:"limit"`
	Skip    *int                `json:"skip" bson:"skip"`
}
