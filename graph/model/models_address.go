package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	Lat      float64                `json:"lat" bson:"lat" form:"osmId"`
	Lon      float64                `json:"lon" bson:"lon" form:"lon"`
	OsmID    string                 `json:"osmId" bson:"osmId" form:"osmId"`
	Address  map[string]interface{} `json:"address" bson:"address" form:"address"`
	DAddress string                 `json:"dAddress" bson:"dAddress" form:"dAddress"`
	Lang     string                 `json:"lang" bson:"lang" form:"lang"`
	Props    map[string]interface{} `json:"props" bson:"props" form:"props"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt" form:"updatedAt"`
}

type AddressFilter struct {
	UserID *string                    `json:"userId,omitempty"`
	OsmID  *string                    `json:"osmId,omitempty"`
	Lon    float64                    `json:"lon" bson:"lon"`
	Lat    float64                    `json:"lat" bson:"lat"`
	Sort   []*ProductFilterSortParams `json:"sort,omitempty"`
	Limit  *int                       `json:"limit,omitempty"`
	Skip   *int                       `json:"skip,omitempty"`
}
