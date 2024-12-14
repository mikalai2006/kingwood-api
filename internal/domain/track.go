package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Track struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`

	Lon float64 `json:"lon" bson:"lon" binding:"required"`
	Lat float64 `json:"lat" bson:"lat" binding:"required"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type TrackInputData struct {
	Lon float64 `json:"lon" bson:"lon" form:"lon"`
	Lat float64 `json:"lat" bson:"lat" form:"lat"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
