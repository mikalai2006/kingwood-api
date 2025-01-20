package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Operation struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	Name   string `json:"name" bson:"name" form:"name"`
	Color  string `json:"color" bson:"color" form:"color"`
	Group  string `json:"group" bson:"group" form:"group"`
	Hidden int    `json:"hidden" bson:"hidden"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OperationInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	Name   string `json:"name" bson:"name" form:"name"`
	Color  string `json:"color" bson:"color" form:"color"`
	Group  string `json:"group" bson:"group" form:"group"`
	Hidden *int   `json:"hidden" bson:"hidden" form:"hidden"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OperationInputData struct {
	Name  string `json:"name" bson:"name" form:"name"`
	Color string `json:"color" bson:"color" form:"color"`
	Group string `json:"group" bson:"group" form:"group"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
