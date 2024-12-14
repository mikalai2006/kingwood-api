package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`
	// ProductID primitive.ObjectID     `json:"productId" bson:"productId"`
	RoomID  primitive.ObjectID     `json:"roomId" bson:"roomId"`
	Status  int                    `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Props   map[string]interface{} `json:"props" bson:"props"`

	Images []string `json:"images" bson:"images"`
	// User User `json:"user,omitempty" bson:"user,omitempty"`
	// Images []MessageImage `json:"images" bson:"images,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type MessageInputMongo struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`
	// ProductID primitive.ObjectID     `json:"productId" bson:"productId"`
	RoomID  primitive.ObjectID     `json:"roomId" bson:"roomId"`
	Status  int                    `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Props   map[string]interface{} `json:"props" bson:"props"`

	Images []string `json:"images" bson:"images"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type MessageImage struct {
	UserID    string `json:"userId" bson:"userId"`
	ServiceID string `json:"serviceId" bson:"serviceId"`
	Service   string `json:"service" bson:"service"`
	Path      string `json:"path" bson:"path"`
	Ext       string `json:"ext" bson:"ext"`
	URL       string `json:"url" bson:"url"`
}

type MessageInput struct {
	// ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID string `json:"userId" bson:"userId"`
	// ProductID primitive.ObjectID     `json:"productId" bson:"productId"`
	RoomID  string                 `json:"roomId" bson:"roomId" form:"roomId" primitive:"true"`
	Status  int                    `json:"status" bson:"status" form:"status"`
	Message string                 `json:"message" bson:"message" form:"message"`
	Props   map[string]interface{} `json:"props" bson:"props" form:"props"`

	Images []string `json:"images" bson:"images"`

	// CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	// UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type MessageFilter struct {
	ID     *primitive.ObjectID `json:"id,omitempty"`
	UserID *primitive.ObjectID `json:"userId,omitempty"`
	// ProductID *primitive.ObjectID        `json:"productId,omitempty"`
	RoomID []*primitive.ObjectID      `json:"roomId" bson:"roomId"`
	Sort   []*ProductFilterSortParams `json:"sort,omitempty"`
	Limit  *int                       `json:"limit,omitempty"`
	Skip   *int                       `json:"skip,omitempty"`
}

type MessageGroupForUser struct {
	UserID    *primitive.ObjectID `json:"userId,omitempty" bson:"userId"`
	ProductID *primitive.ObjectID `json:"productId,omitempty" bson:"productId"`
	// UserProductID primitive.ObjectID  `json:"userProductId" bson:"userProductId"`
	Count   int     `json:"count" bson:"count"`
	Product Product `json:"product" bson:"product"`
}
