package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskWorker struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	ObjectId    primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId     primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	TaskId      primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId    primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	SortOrder   *int64             `json:"sortOrder" bson:"sortOrder"`
	StatusId    primitive.ObjectID `json:"statusId" bson:"statusId"`
	Status      string             `json:"status" bson:"status"`
	From        *time.Time         `json:"from" bson:"from"`
	To          *time.Time         `json:"to" bson:"to"`
	TypeGo      string             `json:"typeGo" bson:"typeGo"`

	Task       Task       `json:"task" bson:"task"`
	TaskStatus TaskStatus `json:"taskStatus" bson:"taskStatus"`
	Order      Order      `json:"order" bson:"order"`
	Worker     User       `json:"worker" bson:"worker"`
	Object     Object     `json:"object" bson:"object"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TaskWorkerInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	ObjectId    primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId     primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	TaskId      primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId    primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	SortOrder   *int64             `json:"sortOrder" bson:"sortOrder"`
	StatusId    primitive.ObjectID `json:"statusId" bson:"statusId" primitive:"true"`
	Status      string             `json:"status" bson:"status"`
	From        time.Time          `json:"from" bson:"from"`
	To          time.Time          `json:"to" bson:"to"`
	TypeGo      string             `json:"typeGo" bson:"typeGo"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TaskWorkerInputData struct {
	ObjectId    string    `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId     string    `json:"orderId" bson:"orderId" primitive:"true"`
	TaskId      string    `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId    string    `json:"workerId" bson:"workerId" primitive:"true"`
	OperationId string    `json:"operationId" bson:"operationId" primitive:"true"`
	SortOrder   *int64    `json:"sortOrder" bson:"sortOrder"`
	StatusId    string    `json:"statusId" bson:"statusId"`
	Status      string    `json:"status" bson:"status"`
	From        time.Time `json:"from" bson:"from"`
	To          time.Time `json:"to" bson:"to"`
	TypeGo      string    `json:"typeGo" bson:"typeGo"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TaskWorkerFilter struct {
	ID          []*string           `json:"id,omitempty"`
	ObjectId    []*string           `json:"objectId,omitempty"`
	OrderId     []*string           `json:"orderId,omitempty"`
	TaskId      []*string           `json:"taskId,omitempty"`
	WorkerId    []*string           `json:"workerId,omitempty"`
	OperationId []*string           `json:"operationId,omitempty"`
	From        *time.Time          `json:"from,omitempty"`
	To          *time.Time          `json:"to,omitempty"`
	Date        *time.Time          `json:"date,omitempty"`
	Sort        []*FilterSortParams `json:"sort,omitempty"`
	Limit       *int                `json:"$limit,omitempty"`
	Skip        *int                `json:"$skip,omitempty"`
}
