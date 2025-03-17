package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArchiveImage struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	ServiceID   string             `json:"serviceId" bson:"serviceId"`
	Service     string             `json:"service" bson:"service"`
	Path        string             `json:"path" bson:"path"`
	Ext         string             `json:"ext" bson:"ext"`
	Title       string             `json:"title" bson:"title"`
	Dir         string             `json:"dir" bson:"dir"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`

	Meta ArchiveMeta `json:"meta" bson:"meta"`
}

type ArchiveImageInput struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	ServiceID   string             `json:"serviceId" bson:"serviceId"`
	Service     string             `json:"service" bson:"service"`
	Path        string             `json:"path" bson:"path"`
	Ext         string             `json:"ext" bson:"ext"`
	Title       string             `json:"title" bson:"title"`
	Dir         string             `json:"dir" bson:"dir"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`

	Meta ArchiveMeta `json:"meta" bson:"meta"`
}
