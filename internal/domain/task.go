package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	OrderId primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	// OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	Name string `json:"name" bson:"name"`
	// WorkerId primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	StartAt  time.Time          `json:"startAt" bson:"startAt"`
	SortOder *int64             `json:"sortOrder" bson:"sortOrder"`
	StatusId primitive.ObjectID `json:"statusId" bson:"statusId"`
	Status   string             `json:"status" bson:"status"`
	Active   *int64             `json:"active" bson:"active" form:"active"`

	AutoCheck *int64       `json:"autoCheck" bson:"autoCheck" form:"autoCheck"`
	Workers   []TaskWorker `json:"-" bson:"workers"`

	From   time.Time `json:"from" bson:"from" form:"from"`
	To     time.Time `json:"to" bson:"to" form:"to"`
	TypeGo string    `json:"typeGo" bson:"typeGo" form:"typeGo"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TaskInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	OrderId primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	Name    string             `json:"name" bson:"name"`
	// OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	// WorkerId primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	StartAt  time.Time          `json:"startAt" bson:"startAt"`
	SortOder *int64             `json:"sortOrder" bson:"sortOrder"`
	StatusId primitive.ObjectID `json:"statusId" bson:"statusId"`
	Status   string             `json:"status" bson:"status"`

	Active    *int64    `json:"active" bson:"active" form:"active"`
	AutoCheck *int64    `json:"autoCheck" bson:"autoCheck" form:"autoCheck"`
	From      time.Time `json:"from" bson:"from" form:"from"`
	To        time.Time `json:"to" bson:"to" form:"to"`
	TypeGo    string    `json:"typeGo" bson:"typeGo" form:"typeGo"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TaskInputData struct {
	ID      string `json:"id" bson:"_id" primitive:"true"`
	OrderId string `json:"orderId" bson:"orderId" primitive:"true"`
	Name    string `json:"name" bson:"name" form:"name"`
	// OperationId string    `json:"operationId" bson:"operationId" primitive:"true"`
	// WorkerId    string    `json:"workerId" bson:"workerId" primitive:"true"`
	// StartAt     time.Time `json:"startAt" bson:"startAt"`
	Status    string `json:"status" bson:"status"`
	SortOder  *int64 `json:"sortOrder" bson:"sortOrder"`
	Active    *int64 `json:"active" bson:"active" form:"active"`
	AutoCheck *int64 `json:"autoCheck" bson:"autoCheck" form:"autoCheck"`
	// Status   string `json:"status" bson:"status"`
	From   time.Time `json:"from" bson:"from" form:"from"`
	To     time.Time `json:"to" bson:"to" form:"to"`
	TypeGo string    `json:"typeGo" bson:"typeGo" form:"typeGo"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
