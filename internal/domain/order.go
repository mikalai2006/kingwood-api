package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	Number        int                `json:"number" bson:"number"`
	Name          string             `json:"name" bson:"name" form:"name"`
	Description   string             `json:"description" bson:"description" form:"description"`
	ConstructorId primitive.ObjectID `json:"constructorId" bson:"constructorId" primitive:"true"`
	ObjectId      primitive.ObjectID `json:"objectId" bson:"objectId" form:"objectId"`
	Term          time.Time          `json:"term" bson:"term" form:"term"`
	Priority      *int64             `json:"priority" bson:"priority" form:"priority"`
	Status        []string           `json:"status" bson:"status" form:"status"`

	// User User `json:"user" bson:"user"`
	Object Object `json:"object" bson:"object"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OrderInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	Number        int64              `json:"number" bson:"number"`
	Name          string             `json:"name" bson:"name" form:"name"`
	Description   string             `json:"description" bson:"description" form:"description"`
	ObjectId      primitive.ObjectID `json:"objectId" bson:"objectId" form:"objectId"`
	ConstructorId primitive.ObjectID `json:"constructorId" bson:"constructorId" primitive:"true"`
	Term          time.Time          `json:"term" bson:"term" form:"term"`
	Priority      *int64             `json:"priority" bson:"priority" form:"priority"`
	Status        []string           `json:"status" bson:"status" form:"status"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OrderInputData struct {
	ID     string   `json:"id" bson:"_id" primitive:"true"`
	Name   string   `json:"name" bson:"name"  form:"name"`
	Status []string `json:"status" bson:"status" form:"status"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
