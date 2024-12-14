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

	// MaxDistance int    `json:"maxDistance" bson:"maxDistance"`
	// Post []Post `json:"post" bson:"post,omitempty" gorm:"-"`

	UserData  User      `json:"-" bson:"userData"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

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

	// VkID     string `json:"vkId" bson:"vk_id" form:"vkId"`
	// AppleID  string `json:"appleId" bson:"apple_id" form:"appleId"`
	// GoogleID string `json:"googleId" bson:"google_id" form:"googleId"`
	// GithubID string `json:"githubId" bson:"github_id" form:"githubId"`

	Verification Verification `json:"verification" bson:"verification"`
	Session      Session      `json:"session" bson:"session"`
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
	// Roles        []string     `json:"roles" bson:"roles"`
	// MaxDistance int `json:"maxDistance" bson:"max_distance"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type DataForClaims struct {
	Roles  []string `json:"roles" bson:"roles"`
	UserID string   `json:"user_id" bson:"user_id"`
	Md     int      `json:"md" bson:"md"`
	UID    string   `json:"uid" bson:"uid"`
}

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}
