package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkHistory struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	ObjectId     primitive.ObjectID     `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId      primitive.ObjectID     `json:"orderId" bson:"orderId" primitive:"true"`
	TaskId       primitive.ObjectID     `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId     primitive.ObjectID     `json:"workerId" bson:"workerId" primitive:"true"`
	OperationId  primitive.ObjectID     `json:"operationId" bson:"operationId" primitive:"true"`
	TaskWorkerId primitive.ObjectID     `json:"taskWorkerId" bson:"taskWorkerId" primitive:"true"`
	Status       int                    `json:"status" bson:"status"`
	Date         time.Time              `json:"date" bson:"date"`
	From         time.Time              `json:"from" bson:"from"`
	To           time.Time              `json:"to" bson:"to"`
	Oklad        *int64                 `json:"oklad" bson:"oklad"`
	Total        *int64                 `json:"total" bson:"total"`
	TotalTime    *int64                 `json:"totalTime" bson:"totalTime"`
	Props        map[string]interface{} `json:"props" bson:"props" form:"props"`

	// Task Task `json:"task" bson:"task"`
	// TaskStatus TaskStatus `json:"taskStatus" bson:"taskStatus"`
	Order  Order  `json:"order" bson:"order"`
	Worker User   `json:"worker" bson:"worker"`
	Object Object `json:"object" bson:"object"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type WorkHistoryInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"userId"`

	ObjectId     *primitive.ObjectID `json:"objectId" bson:"objectId" primitive:"true"`
	OrderId      *primitive.ObjectID `json:"orderId" bson:"orderId" primitive:"true"`
	TaskId       *primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkerId     primitive.ObjectID  `json:"workerId" bson:"workerId" primitive:"true"`
	OperationId  *primitive.ObjectID `json:"operationId" bson:"operationId" primitive:"true"`
	TaskWorkerId *primitive.ObjectID `json:"taskWorkerId" bson:"taskWorkerId" primitive:"true"`

	Status    *int                   `json:"status" bson:"status"`
	Date      time.Time              `json:"date" bson:"date"`
	From      time.Time              `json:"from" bson:"from"`
	To        time.Time              `json:"to" bson:"to"`
	Oklad     *int64                 `json:"oklad" bson:"oklad"`
	Total     *int64                 `json:"total" bson:"total"`
	TotalTime *int64                 `json:"totalTime" bson:"totalTime"`
	Props     map[string]interface{} `json:"props" bson:"props" form:"props"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type WorkHistoryFilter struct {
	ID []string `json:"id,omitempty"`
	// OrderId    []string            `json:"orderId,omitempty"`
	WorkerId     []string            `json:"workerId,omitempty"`
	TaskWorkerId []string            `json:"taskWorkerId,omitempty"`
	TaskId       []string            `json:"taskId,omitempty"`
	OrderId      []string            `json:"orderId,omitempty"`
	Status       *int                `json:"status,omitempty"`
	From         time.Time           `json:"from,omitempty"`
	To           time.Time           `json:"to,omitempty"`
	Date         time.Time           `json:"date,omitempty"`
	Sort         []*FilterSortParams `json:"sort,omitempty"`
	Limit        *int                `json:"$limit,omitempty"`
	Skip         *int                `json:"$skip,omitempty"`
}

type WorkHistoryStatByOrderOperation struct {
	OperationId primitive.ObjectID `json:"operationId" bson:"operationId"`
	Count       int64              `json:"count" bson:"count"`
	Total       int64              `json:"total" bson:"total"`
	// Operation Operation          `json:"operation" bson:"operation"`
}

type WorkHistoryStatByOrder struct {
	ID         primitive.ObjectID                `json:"workerId" bson:"_id"`
	Count      int64                             `json:"count" bson:"count"`
	Total      int64                             `json:"total" bson:"total"`
	Worker     User                              `json:"worker" bson:"worker"`
	Operations []WorkHistoryStatByOrderOperation `json:"operations" bson:"operations"`
}

// type WorkHistoryStatByMonthOrder struct {
// 	OperationId primitive.ObjectID `json:"operationId" bson:"operationId"`
// 	Count       int64              `json:"count" bson:"count"`
// 	Total       int64              `json:"total" bson:"total"`
// 	// Operation Operation          `json:"operation" bson:"operation"`
// }

type WorkHistoryStatByMonth struct {
	ID      primitive.ObjectID `json:"orderId" bson:"_id"`
	Count   int64              `json:"count" bson:"count"`
	Total   int64              `json:"total" bson:"total"`
	TotalMs int64              `json:"totalTime" bson:"totalMs"`
	Order   Order              `json:"order" bson:"order"`
	Workers []string           `json:"workers" bson:"workers"`
	// Operations []WorkHistoryStatByOrderOperation `json:"operations" bson:"operations"`
}
