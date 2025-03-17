package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArchiveTask struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID  `json:"userId" bson:"userId" primitive:"true"`
	ObjectId    primitive.ObjectID  `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId     primitive.ObjectID  `json:"orderId" bson:"orderId" primitive:"true"`
	OperationId primitive.ObjectID  `json:"operationId" bson:"operationId" primitive:"true"`
	Name        string              `json:"name" bson:"name"`
	StartAt     time.Time           `json:"startAt" bson:"startAt"`
	SortOrder   *int64              `json:"sortOrder" bson:"sortOrder"`
	StatusId    primitive.ObjectID  `json:"statusId" bson:"statusId"`
	Status      string              `json:"status" bson:"status"`
	Active      *int64              `json:"active" bson:"active" form:"active"`
	AutoCheck   *int64              `json:"autoCheck" bson:"autoCheck" form:"autoCheck"`
	Workers     []ArchiveTaskWorker `json:"workers" bson:"workers"`
	Object      Object              `json:"object" bson:"object"`
	Operation   Operation           `json:"operation" bson:"operation"`
	Order       Order               `json:"order" bson:"order"`
	From        *time.Time          `json:"from" bson:"from" form:"from"`
	To          *time.Time          `json:"to" bson:"to" form:"to"`
	TypeGo      string              `json:"typeGo" bson:"typeGo" form:"typeGo"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`

	Meta ArchiveMeta `json:"meta" bson:"meta"`
}

type ArchiveTaskInput struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	ObjectId    primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId     primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	Name        string             `json:"name" bson:"name"`
	OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	StartAt     time.Time          `json:"startAt" bson:"startAt"`
	SortOrder   *int64             `json:"sortOrder" bson:"sortOrder"`
	StatusId    primitive.ObjectID `json:"statusId" bson:"statusId"`
	Status      string             `json:"status" bson:"status"`
	Active      *int64             `json:"active" bson:"active" form:"active"`
	AutoCheck   *int64             `json:"autoCheck" bson:"autoCheck" form:"autoCheck"`
	From        time.Time          `json:"from" bson:"from" form:"from"`
	To          time.Time          `json:"to" bson:"to" form:"to"`
	TypeGo      string             `json:"typeGo" bson:"typeGo" form:"typeGo"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`

	Meta ArchiveMeta `json:"meta" bson:"meta"`
}

type ArchiveTaskFilter struct {
	ID          []string            `json:"id,omitempty"`
	ObjectId    []string            `json:"objectId,omitempty"`
	OrderId     []string            `json:"orderId,omitempty"`
	OperationId []string            `json:"operationId,omitempty"`
	Status      []string            `json:"status,omitempty"`
	Name        string              `json:"name,omitempty"`
	Sort        []*FilterSortParams `json:"sort,omitempty"`
	Limit       *int                `json:"$limit,omitempty"`
	Skip        *int                `json:"$skip,omitempty"`
}
