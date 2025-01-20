package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkHistory struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	WorkTimeId  primitive.ObjectID `json:"workTimeId" bson:"workTimeId" primitive:"true"`
	ObjectId    primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId     primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	TaskId      primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId    primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	Status      int                `json:"status" bson:"status"`
	From        time.Time          `json:"from" bson:"from"`
	To          time.Time          `json:"to" bson:"to"`
	Oklad       *int64             `json:"oklad" bson:"oklad"`

	// Task       Task       `json:"task" bson:"task"`
	// TaskStatus TaskStatus `json:"taskStatus" bson:"taskStatus"`
	// Order      Order      `json:"order" bson:"order"`
	// Worker     User       `json:"worker" bson:"worker"`
	// Object     Object     `json:"object" bson:"object"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type WorkHistoryInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	WorkTimeId  primitive.ObjectID `json:"workTimeId" bson:"workTimeId" primitive:"true"`
	ObjectId    primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId     primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	TaskId      primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId    primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	OperationId primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	Status      *int               `json:"status" bson:"status"`
	From        time.Time          `json:"from" bson:"from"`
	To          time.Time          `json:"to" bson:"to"`
	Oklad       *int64             `json:"oklad" bson:"oklad"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type WorkHistoryFilter struct {
	ID         []string `json:"id,omitempty"`
	WorkTimeId []string `json:"workTimeId,omitempty"`
	// OrderId    []string            `json:"orderId,omitempty"`
	WorkerId []string            `json:"workerId,omitempty"`
	TaskId   []string            `json:"taskId,omitempty"`
	Status   *int                `json:"status,omitempty"`
	Sort     []*FilterSortParams `json:"sort,omitempty"`
	Limit    *int                `json:"$limit,omitempty"`
	Skip     *int                `json:"$skip,omitempty"`
}
