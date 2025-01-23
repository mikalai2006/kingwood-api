package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	// swagger:ignore
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login string             `json:"login" bson:"login" binding:"required"`
	Email string             `json:"email" bson:"email"`
	// Phone        string             `json:"phone" bson:"phone"`
	Password string `json:"password" bson:"password" binding:"required"`
	Strategy string `json:"strategy" bson:"strategy"`
	// VkID         string       `json:"vkId" bson:"vkId"`
	// AppleID      string       `json:"appleId" bson:"appleId"`
	// GoogleID     string       `json:"googleId" bson:"googleId"`
	// GithubID     string       `json:"githubId" bson:"githubId"`
	Verification Verification `json:"verification" bson:"verification"`
	Session      Session      `json:"session" bson:"session"`
	PushToken    string       `json:"pushToken" bson:"pushToken"`

	// MaxDistance int    `json:"maxDistance" bson:"maxDistance"`
	// Post []Post `json:"post" bson:"post,omitempty" gorm:"-"`

	Role      Role      `json:"-" bson:"role"`
	User      User      `json:"-" bson:"user"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
type AuthPrivate struct {
	// swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login        string             `json:"login" bson:"login"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Strategy     string             `json:"strategy" bson:"strategy"`
	Verification Verification       `json:"verification" bson:"verification"`
	Session      Session            `json:"session" bson:"session"`
	PushToken    string             `json:"pushToken" bson:"pushToken"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// type AuthResetPassword struct {
// 	UserId   primitive.ObjectID `json:"userId" bson:"userId" form:"userId"`
// }

type AuthInput struct {
	// UserId primitive.ObjectID `json:"userId" bson:"userId" form:"userId" gorm:"-"`
	RoleId   primitive.ObjectID `json:"roleId" bson:"roleId" form:"roleId"`
	PostId   primitive.ObjectID `json:"postId" bson:"postId" form:"postId"`
	TypeWork []string           `json:"typeWork" bson:"typeWork" form:"typeWork"`

	Login string `json:"login" bson:"login" form:"login"`
	Email string `json:"email" bson:"email" form:"email"`

	Name  string `json:"name" bson:"name" form:"name"`
	Phone string `json:"phone" bson:"phone"`
	// Phone string `json:"phone" bson:"phone" form:"phone"`
	// Post     []Post  `json:"post" bson:"post" form:"post"`
	Password string `json:"password" bson:"password" form:"password"`
	Strategy string `json:"strategy" bson:"strategy"`
	TypePay  *int64 `json:"typePay" bson:"typePay" form:"typePay"`
	Oklad    *int64 `json:"oklad" bson:"oklad" form:"oklad"`
	Birthday string `json:"birthday" bson:"birthday" form:"birthday"`

	// VkID     string `json:"vkId" bson:"vk_id" form:"vkId"`
	// AppleID  string `json:"appleId" bson:"apple_id" form:"appleId"`
	// GoogleID string `json:"googleId" bson:"google_id" form:"googleId"`
	// GithubID string `json:"githubId" bson:"github_id" form:"githubId"`

	Verification Verification `json:"verification" bson:"verification"`
	Session      Session      `json:"session" bson:"session"`
	PushToken    string       `json:"pushToken" bson:"pushToken"`
	// Roles        []string     `json:"roles" bson:"roles"`
	// MaxDistance int `json:"maxDistance" bson:"max_distance"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AuthInputMongo struct {
	// UserId primitive.ObjectID `json:"userId" bson:"userId" form:"userId" gorm:"-"`
	// RoleId primitive.ObjectID `json:"roleId" bson:"roleId" form:"roleId"`
	// PostId primitive.ObjectID `json:"postId" bson:"postId" form:"postId"`

	Login string `json:"login" bson:"login" form:"login"`
	Email string `json:"email" bson:"email" form:"email"`

	// Phone string `json:"phone" bson:"phone" form:"phone"`
	// Post     []Post  `json:"post" bson:"post" form:"post"`
	Password string `json:"password" bson:"password" form:"password"`
	Strategy string `json:"strategy" bson:"strategy"`
	// VkID     string `json:"vkId" bson:"vk_id" form:"vkId"`
	// AppleID  string `json:"appleId" bson:"apple_id" form:"appleId"`
	// GoogleID string `json:"googleId" bson:"google_id" form:"googleId"`
	// GithubID string `json:"githubId" bson:"github_id" form:"githubId"`

	Verification Verification `json:"verification" bson:"verification"`
	Session      Session      `json:"session" bson:"session"`
	PushToken    string       `json:"pushToken" bson:"pushToken"`
	// Roles        []string     `json:"roles" bson:"roles"`
	// MaxDistance int `json:"maxDistance" bson:"max_distance"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type DataForClaims struct {
	// Roles  []string `json:"roles" bson:"roles"`
	UserID string `json:"userId" bson:"userId"`
	// Md     int      `json:"md" bson:"md"`
	UID string `json:"uid" bson:"uid"`
}

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}

type AuthPublicData struct {
	Login     string `json:"login" bson:"login"`
	PushToken string `json:"pushToken" bson:"pushToken"`
}

type ResetPassword struct {
	Password string `json:"password" bson:"password"`
}
