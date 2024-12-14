package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Action struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID     `json:"userId" bson:"user_id"`
	Service     string                 `json:"service" bson:"service"`
	ServiceID   primitive.ObjectID     `json:"serviceId" bson:"service_id"`
	Type        int64                  `json:"type" bson:"type"`
	Description string                 `json:"description" bson:"description"`
	Props       map[string]interface{} `json:"props" bson:"props"`
	Status      int64                  `json:"status" bson:"status"`
	CreatedAt   time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updated_at"`
}

type ActionInput struct {
	UserID      string                 `json:"userId" bson:"user_id"`
	Service     string                 `json:"service" bson:"service"`
	ServiceID   string                 `json:"serviceId" bson:"service_id"`
	Type        int64                  `json:"type" bson:"type"`
	Description string                 `json:"description" bson:"description"`
	Props       map[string]interface{} `json:"props" bson:"props"`
	Status      int64                  `json:"status" bson:"status"`
}
