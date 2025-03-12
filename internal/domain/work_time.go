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
	Status   int                `json:"status" bson:"status"`
	Date     time.Time          `json:"date" bson:"date"`
	From     time.Time          `json:"from" bson:"from"`
	To       time.Time          `json:"to" bson:"to"`
	Oklad    *int64             `json:"oklad" bson:"oklad"`
	Total    *int64             `json:"total" bson:"total"`

	Props map[string]interface{} `json:"props" bson:"props" form:"props"`
	// Task       Task       `json:"task" bson:"task"`
	// TaskStatus TaskStatus `json:"taskStatus" bson:"taskStatus"`
	// Order      Order      `json:"order" bson:"order"`
	WorHistory []WorkHistory `json:"workHistory" bson:"workHistory"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type WorkTimeInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	// OrderId  primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	// TaskId   primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	Status   *int               `json:"status" bson:"status"`
	Oklad    *int64             `json:"oklad" bson:"oklad"`
	Total    *int64             `json:"total" bson:"total"`

	Props map[string]interface{} `json:"props" bson:"props" form:"props"`

	Date time.Time `json:"date" bson:"date"`
	From time.Time `json:"from" bson:"from"`
	To   time.Time `json:"to" bson:"to"`

	OrderId string `json:"orderId" bson:"orderId" primitive:"true"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// type WorkTimeInputData struct {
// 	// OrderId  string    `json:"orderId" bson:"orderId" primitive:"true"`
// 	// TaskId   string    `json:"taskId" bson:"taskId" primitive:"true"`
// 	WorkerId string `json:"workerId" bson:"workerId" primitive:"true"`
// 	Status   int `json:"status" bson:"status"`

// 	Date time.Time `json:"date" bson:"date"`
// 	From time.Time `json:"from" bson:"from"`
// 	To   time.Time `json:"to" bson:"to"`

// 	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
// }

type WorkTimeFilter struct {
	ID       []string            `json:"id,omitempty"`
	WorkerId []string            `json:"workerId,omitempty"`
	From     time.Time           `json:"from,omitempty"`
	To       time.Time           `json:"to,omitempty"`
	Date     time.Time           `json:"date,omitempty"`
	Status   *int                `json:"status,omitempty"`
	Sort     []*FilterSortParams `json:"sort,omitempty"`
	Limit    *int                `json:"$limit,omitempty"`
	Skip     *int                `json:"$skip,omitempty"`
}
