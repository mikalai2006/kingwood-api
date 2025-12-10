package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStat struct {
}

type User struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty" primitive:"true"`

	Name  string `json:"name" bson:"name" form:"name"`
	Phone string `json:"phone" bson:"phone"`
	// Avatar string `json:"avatar" bson:"avatar"`
	Online   *bool   `json:"online" bson:"online" form:"online"`
	Hidden   int     `json:"hidden" bson:"hidden" form:"hidden"`
	Archive  int     `json:"archive" bson:"archive" form:"archive"`
	Birthday *string `json:"birthday" bson:"birthday" form:"birthday"`
	// Post   []int  `json:"post" bson:"post"`
	// Verify   bool   `json:"verify" bson:"verify"`
	// Login    string `json:"login" bson:"login" form:"login"`
	// Location GeoLocation `json:"location" bson:"location" form:"location"`
	// UserStat UserStat `json:"userStat" bson:"user_stat"`

	// Md     int      `json:"md" bson:"md"`
	// Bal    int      `json:"bal" bson:"bal"`
	// Role   Role     `json:"role" bson:"role"`

	RoleId primitive.ObjectID `json:"roleId" bson:"roleId" form:"roleId" primitive:"true"`
	// RolesId  []primitive.ObjectID   `json:"rolesId" bson:"rolesId" form:"rolesId" primitive:"true"`
	PostId   primitive.ObjectID     `json:"postId" bson:"postId" form:"postId"`
	TypeWork []string               `json:"typeWork" bson:"typeWork"`
	TypePay  *int64                 `json:"typePay" bson:"typePay"`
	Oklad    *int64                 `json:"oklad" bson:"oklad"`
	MaxTime  *int64                 `json:"maxTime" bson:"maxTime"`
	Props    map[string]interface{} `json:"props" bson:"props" form:"props"`
	Blocked  *int                   `json:"blocked" bson:"blocked" form:"blocked"`

	TaskWorkers []TaskWorker `json:"taskWorkers" bson:"taskWorkers" form:"taskWorkers"`
	// Workes     *int64             `json:"workes" bson:"workes"`

	PostObject   Post          `json:"postObject" bson:"postObject"`
	RoleObject   Role          `json:"roleObject" bson:"roleObject"`
	Roles        []Role        `json:"roles" bson:"roles"`
	Images       []Image       `json:"images,omitempty" bson:"images,omitempty"`
	IsWork       int           `json:"isWork" bson:"isWork"`
	WorkHistorys []WorkHistory `json:"workHistorys" bson:"workHistorys"`
	// Post   []string `json:"post" bson:"post"`
	Auth        AuthPublicData `json:"auth,omitempty" bson:"auth,omitempty"`
	AuthPrivate AuthPrivate    `json:"-" bson:"authPrivate,omitempty"`
	Dops        []UserDoplata  `json:"dops" bson:"dops"`

	LastTime  time.Time `json:"lastTime" bson:"lastTime"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UserFlat struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty" primitive:"true"`

	Name     string  `json:"name" bson:"name" form:"name"`
	Phone    string  `json:"phone" bson:"phone"`
	Online   *bool   `json:"online" bson:"online" form:"online"`
	Hidden   int     `json:"hidden" bson:"hidden" form:"hidden"`
	Archive  int     `json:"archive" bson:"archive" form:"archive"`
	Birthday *string `json:"birthday" bson:"birthday" form:"birthday"`

	RoleId primitive.ObjectID `json:"roleId" bson:"roleId" form:"roleId" primitive:"true"`
	// RolesId  []primitive.ObjectID   `json:"rolesId" bson:"rolesId" form:"rolesId" primitive:"true"`
	PostId   primitive.ObjectID     `json:"postId" bson:"postId" form:"postId"`
	TypeWork []string               `json:"typeWork" bson:"typeWork"`
	TypePay  *int64                 `json:"typePay" bson:"typePay"`
	Oklad    *int64                 `json:"oklad" bson:"oklad"`
	MaxTime  *int64                 `json:"maxTime" bson:"maxTime"`
	Props    map[string]interface{} `json:"props" bson:"props" form:"props"`
	Blocked  *int                   `json:"blocked" bson:"blocked" form:"blocked"`

	// // Workes     *int64             `json:"workes" bson:"workes"`

	// PostObject   Post          `json:"postObject" bson:"postObject"`
	// RoleObject   Role          `json:"roleObject" bson:"roleObject"`
	// Roles        []Role        `json:"roles" bson:"roles"`
	// Images       []Image       `json:"images,omitempty" bson:"images,omitempty"`
	// IsWork       int           `json:"isWork" bson:"isWork"`
	// WorkHistorys []WorkHistory `json:"workHistorys" bson:"workHistorys"`
	// // Post   []string `json:"post" bson:"post"`
	// Auth        AuthPublicData `json:"auth,omitempty" bson:"auth,omitempty"`
	// AuthPrivate AuthPrivate    `json:"-" bson:"authPrivate,omitempty"`
	// Dops        []UserDoplata  `json:"dops" bson:"dops"`

	// LastTime  time.Time `json:"lastTime" bson:"lastTime"`
	// CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	// UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UserInput struct {
	// ID     primitive.ObjectID `json:"id" bson:"_id" form:"id" primitive:"true"`
	// UserID primitive.ObjectID `json:"userId" bson:"userId" form:"userId"`
	Name string `json:"name" bson:"name" form:"name"`
	// Login    string `json:"login" bson:"login" form:"login"`
	Phone string `json:"phone" bson:"phone" form:"phone"`
	// Avatar string `json:"avatar" bson:"avatar" form:"avatar"`
	Hidden  *int `json:"hidden" bson:"hidden" form:"hidden"`
	Archive *int `json:"archive" bson:"archive" form:"archive"`
	// Post   []int  `json:"post" bson:"post" form:"post"`
	RoleId string `json:"roleId" bson:"roleId" form:"roleId"`
	// RolesId  []string               `json:"rolesId" bson:"rolesId" form:"rolesId"`
	PostId   string                 `json:"postId" bson:"postId" form:"postId"`
	TypeWork []string               `json:"typeWork" bson:"typeWork" form:"typeWork"`
	Birthday *string                `json:"birthday" bson:"birthday" form:"birthday"`
	Online   *bool                  `json:"online" bson:"online" form:"online"`
	Props    map[string]interface{} `json:"props" bson:"props" form:"props"`
	Blocked  *int                   `json:"blocked" bson:"blocked" form:"blocked"`

	TypePay *int64        `json:"typePay" bson:"typePay" form:"typePay"`
	Oklad   *int64        `json:"oklad" bson:"oklad" form:"oklad"`
	MaxTime *int64        `json:"maxTime" bson:"maxTime" form:"maxTime"`
	Dops    []UserDoplata `json:"dops" bson:"dops" form:"dops"`
	// Workes  *int64 `json:"workes" bson:"workes" form:"workes"`

	LastTime  time.Time `json:"lastTime" bson:"lastTime"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UserInputMongo struct {
	// ID     primitive.ObjectID `json:"id" bson:"_id" form:"id" primitive:"true"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" form:"userId"`
	Name   string             `json:"name" bson:"name" form:"name"`
	// Login    string `json:"login" bson:"login" form:"login"`
	Phone string `json:"phone" bson:"phone" form:"phone"`
	// Avatar string `json:"avatar" bson:"avatar" form:"avatar"`
	Hidden  int `json:"hidden" bson:"hidden" form:"hidden"`
	Archive int `json:"archive" bson:"archive" form:"archive"`
	// Post   []int  `json:"post" bson:"post" form:"post"`
	RoleId primitive.ObjectID `json:"roleId" bson:"roleId" form:"roleId"`
	// RolesId  []primitive.ObjectID `json:"rolesId" bson:"rolesId" form:"rolesId"`
	PostId   primitive.ObjectID `json:"postId" bson:"postId" form:"postId"`
	TypeWork []string           `json:"typeWork" bson:"typeWork"`
	Birthday *string            `json:"birthday" bson:"birthday" form:"birthday"`

	TypePay *int64                 `json:"typePay" bson:"typePay" form:"typePay"`
	Oklad   *int64                 `json:"oklad" bson:"oklad" form:"oklad"`
	MaxTime *int64                 `json:"maxTime" bson:"maxTime" form:"maxTime"`
	Props   map[string]interface{} `json:"props" bson:"props" form:"props"`
	Blocked *int                   `json:"blocked" bson:"blocked" form:"blocked"`
	// Workes  *int64 `json:"workes" bson:"workes" form:"workes"`

	LastTime  time.Time `json:"lastTime" bson:"lastTime"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UserDoplata struct {
	ID       int   `json:"id" bson:"id" form:"id"`
	Days     []int `json:"days" bson:"days" form:"days"`
	MinHours int   `json:"minHours" bson:"minHours" form:"minHours"`
	Doplata  int   `json:"doplata" bson:"doplata" form:"doplata"`
}

type UserFilter struct {
	ID     []string `json:"id,omitempty"`
	UserId []string `json:"userId,omitempty"`
	RoleId []string `json:"roleId,omitempty"`
	// RolesId []string            `json:"rolesId,omitempty"`
	Hidden  *int                `json:"hidden,omitempty"`
	Blocked *int                `json:"blocked,omitempty"`
	Archive *int                `json:"archive,omitempty"`
	Sort    []*FilterSortParams `json:"sort,omitempty"`
	Limit   *int                `json:"$limit,omitempty"`
	Skip    *int                `json:"$skip,omitempty"`
}
