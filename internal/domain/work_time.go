package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkTime struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	// OrderId  primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	// TaskId   primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	Status   string             `json:"status" bson:"status"`
	Date     time.Time          `json:"date" bson:"date"`
	From     time.Time          `json:"from" bson:"from"`
	To       time.Time          `json:"to" bson:"to"`

	// Task       Task       `json:"task" bson:"task"`
	// TaskStatus TaskStatus `json:"taskStatus" bson:"taskStatus"`
	// Order      Order      `json:"order" bson:"order"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type WorkTimeInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	// OrderId  primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	// TaskId   primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	Status   string             `json:"status" bson:"status"`

	Date time.Time `json:"date" bson:"date"`
	From time.Time `json:"from" bson:"from"`
	To   time.Time `json:"to" bson:"to"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type WorkTimeInputData struct {
	// OrderId  string    `json:"orderId" bson:"orderId" primitive:"true"`
	// TaskId   string    `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId string `json:"workerId" bson:"workerId" primitive:"true"`
	Status   string `json:"status" bson:"status"`

	Date time.Time `json:"date" bson:"date"`
	From time.Time `json:"from" bson:"from"`
	To   time.Time `json:"to" bson:"to"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type WorkTimeFilter struct {
	ID       []*string           `json:"id,omitempty"`
	WorkerId []*string           `json:"workerId,omitempty"`
	Sort     []*FilterSortParams `json:"sort,omitempty"`
	Limit    *int                `json:"$limit,omitempty"`
	Skip     *int                `json:"$skip,omitempty"`
}
