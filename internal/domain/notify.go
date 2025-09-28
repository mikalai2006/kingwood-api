package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notify struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID     `json:"userId" bson:"userId"`
	UserTo     primitive.ObjectID     `json:"userTo" bson:"userTo"`
	Status     int                    `json:"status" bson:"status"`
	Title      string                 `json:"title" bson:"title"`
	Message    string                 `json:"message" bson:"message"`
	Link       string                 `json:"link" bson:"link"`
	LinkOption map[string]interface{} `json:"linkOption" bson:"linkOption"`
	Props      map[string]interface{} `json:"props" bson:"props"`

	Images    []string `json:"images" bson:"images"`
	User      User     `json:"user" bson:"user"`
	Recepient User     `json:"recepient" bson:"recepient"`

	ReadAt    time.Time `json:"readAt" bson:"readAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type NotifyInputMongo struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID     `json:"userId" bson:"userId"`
	UserTo     primitive.ObjectID     `json:"userTo" bson:"userTo"`
	Status     int                    `json:"status" bson:"status"`
	Title      string                 `json:"title" bson:"title"`
	Message    string                 `json:"message" bson:"message"`
	Link       string                 `json:"link" bson:"link"`
	LinkOption map[string]interface{} `json:"linkOption" bson:"linkOption"`
	Props      map[string]interface{} `json:"props" bson:"props"`

	Images []string `json:"images" bson:"images"`

	ReadAt    time.Time `json:"readAt" bson:"readAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type NotifyImage struct {
	UserID    string `json:"userId" bson:"userId"`
	ServiceID string `json:"serviceId" bson:"serviceId"`
	Service   string `json:"service" bson:"service"`
	Path      string `json:"path" bson:"path"`
	Ext       string `json:"ext" bson:"ext"`
	URL       string `json:"url" bson:"url"`
}

type NotifyInput struct {
	UserID     string                 `json:"userId" bson:"userId" primitive:"true"`
	UserTo     string                 `json:"userTo" bson:"userTo" form:"userTo" primitive:"true"`
	Status     *int                   `json:"status" bson:"status" form:"status"`
	Title      string                 `json:"title" bson:"title" form:"title"`
	Message    string                 `json:"message" bson:"message" form:"message"`
	Link       string                 `json:"link" bson:"link" form:"link"`
	LinkOption map[string]interface{} `json:"linkOption" bson:"linkOption" form:"linkOption"`
	Props      map[string]interface{} `json:"props" bson:"props" form:"props"`

	Images []string `json:"images" bson:"images"`
}

type NotifyFilter struct {
	ID     []*string           `json:"id,omitempty"`
	UserID []*string           `json:"userId,omitempty"`
	UserTo []*string           `json:"userTo,omitempty"`
	Status *int                `json:"status,omitempty"`
	Sort   []*FilterSortParams `json:"sort,omitempty"`
	Limit  *int                `json:"$limit,omitempty"`
	Skip   *int                `json:"$skip,omitempty"`
}

type NotifyListQuery struct {
	ID []*string `json:"id,omitempty"`
}

type ResultFacetNotify struct {
	Metadata []ResultMetadata `json:"metadata" bson:"metadata"`
	Data     []Notify         `json:"data" bson:"data"`
}
