package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pay struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	Name     string `json:"name" bson:"name"`
	EveryDay *int64 `json:"everyDay" bson:"everyDay"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type PayInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	Name     string `json:"name" bson:"name"`
	EveryDay *int64 `json:"everyDay" bson:"everyDay"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type PayInputData struct {
	Name     string `json:"name" bson:"name"`
	EveryDay *int64 `json:"everyDay" bson:"everyDay"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
