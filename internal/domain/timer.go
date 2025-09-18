package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimerShedule struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	IDTimer   string             `json:"idTimer" bson:"idTimer"`
	ExecuteAt time.Time          `json:"executeAt" bson:"executeAt"`
	IsRunning int                `json:"isRunning" bson:"isRunning"`

	WorkerId      primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	TaskWorkerId  primitive.ObjectID `json:"taskWorkerId" bson:"taskWorkerId" primitive:"true"`
	TaskId        primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkHistoryId primitive.ObjectID `json:"workHistoryId" bson:"workHistoryId" primitive:"true"`
	Worker        User               `json:"worker" bson:"worker"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TimerSheduleInput struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	IDTimer   string             `json:"idTimer" bson:"idTimer"`
	ExecuteAt time.Time          `json:"executeAt" bson:"executeAt"`
	IsRunning *int               `json:"isRunning" bson:"isRunning"`

	WorkerId      primitive.ObjectID `json:"workerId" bson:"workerId" primitive:"true"`
	TaskWorkerId  primitive.ObjectID `json:"taskWorkerId" bson:"taskWorkerId" primitive:"true"`
	TaskId        primitive.ObjectID `json:"taskId" bson:"taskId" primitive:"true"`
	WorkHistoryId primitive.ObjectID `json:"workHistoryId" bson:"workHistoryId" primitive:"true"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TimerSheduleFilter struct {
	ID            []string            `json:"id,omitempty"`
	IDTimer       []string            `json:"idTimer" bson:"idTimer"`
	WorkerId      []string            `json:"workerId,omitempty"`
	TaskId        []string            `json:"taskId,omitempty"`
	TaskWorkerId  []string            `json:"taskWorkerId,omitempty"`
	WorkHistoryId []string            `json:"workHistoryId,omitempty"`
	ExecuteAt     time.Time           `json:"executeAt,omitempty"`
	IsRunning     *int                `json:"isRunning,omitempty"`
	Sort          []*FilterSortParams `json:"sort,omitempty"`
	Limit         *int                `json:"$limit,omitempty"`
	Skip          *int                `json:"$skip,omitempty"`
}
