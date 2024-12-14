package domain

import (
	"time"
)

type StatInfo struct {
	// ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string    `json:"name" bson:"name"`
	CCode         string    `json:"ccode" bson:"ccode"`
	Path          string    `json:"path" bson:"path"`
	Size          float64   `json:"size" bson:"size"`
	Count         float64   `json:"count" bson:"count"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt" bson:"lastUpdatedAt"`
	// UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
}

type NodeFileItem struct {
	ID string `json:"id" bson:"_id,omitempty"`
	// AmenityID string `json:"amenityId" bson:"amenity_id"`
	Type   string  `json:"type" bson:"type"`
	Lon    float64 `json:"lon" bson:"lon"`
	Lat    float64 `json:"lat" bson:"lat"`
	Name   string  `json:"name" bson:"name"`
	CCode  string  `json:"ccode" bson:"ccode"`
	UserID string  `json:"userId" bson:"user_id"`

	Data []NodeFileItemDataItem `json:"data" bson:"data,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	// UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
type NodeFileItemDataItem struct {
	ID        string                   `json:"id" bson:"_id,omitempty"`
	NodeID    string                   `json:"nodeId" bson:"node_id"`
	TagID     string                   `json:"tagId" bson:"tag_id"`
	TagoptID  string                   `json:"tagoptId" bson:"tagopt_id"`
	Data      NodeFileItemDataItemData `json:"data" bson:"data"`
	CreatedAt time.Time                `json:"createdAt" bson:"created_at"`
	// UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type NodeFileItemDataItemData struct {
	Value interface{} `json:"value" bson:"value"`
}
