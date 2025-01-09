package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskMontaj struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	ObjectId    primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
	OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	Name        string             `json:"name" bson:"name"`
	SortOder    *int64             `json:"sortOrder" bson:"sortOrder"`
	StatusId    primitive.ObjectID `json:"statusId" bson:"statusId"`
	Status      string             `json:"status" bson:"status"`
	// Workers  []TaskWorker       `json:"-" bson:"workers"`
	TypeGo string     `json:"typeGo" bson:"typeGo" form:"typeGo"`
	From   *time.Time `json:"from" bson:"from" form:"from"`
	To     time.Time  `json:"to" bson:"to" form:"to"`

	Object Object `json:"object" bson:"object"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TaskMontajInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	ObjectId    primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
	OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	Name        string             `json:"name" bson:"name"`
	SortOder    *int64             `json:"sortOrder" bson:"sortOrder"`
	StatusId    primitive.ObjectID `json:"statusId" bson:"statusId"`
	Status      string             `json:"status" bson:"status"`
	TypeGo      string             `json:"typeGo" bson:"typeGo" form:"typeGo"`
	From        time.Time          `json:"from" bson:"from" form:"from"`
	To          time.Time          `json:"to" bson:"to" form:"to"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// type TaskMontajInputData struct {
// 	ID string `json:"id" bson:"_id" primitive:"true"`
// 	ObjectId primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
// 	Name     string             `json:"name" bson:"name"`
// 	SortOder *int64             `json:"sortOrder" bson:"sortOrder"`
// 	StatusId primitive.ObjectID `json:"statusId" bson:"statusId"`
// 	Status   string             `json:"status" bson:"status"`
// 	TypeGo   string             `json:"typeGo" bson:"typeGo" form:"typeGo"`
// 	From     time.Time          `json:"from" bson:"from" form:"from"`
// 	To       time.Time          `json:"to" bson:"to" form:"to"`
// 	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
// }

type TaskMontajFilter struct {
	ID    []*string           `json:"id,omitempty"`
	Name  *string             `json:"name,omitempty"`
	Sort  []*FilterSortParams `json:"sort,omitempty"`
	Limit *int                `json:"$limit,omitempty"`
	Skip  *int                `json:"$skip,omitempty"`
}
