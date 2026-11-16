package models

import (
	"time"
)

// user db model and validation
type User struct {
	ID         int          `gorm:"primary_key"`
	UserName   string       `json:"username" gorm:"unique"  binding:"required"`
	FullName   string       `json:"fullname"  binding:"required"`
	Email      string       `json:"email" gorm:"unique"  binding:"required,email"`
	Password   string       `json:"password" binding:"required,min=6"`
	Status     bool         `gorm:"default:false"`
	GroupChats []*GroupChat `gorm:"many2many:user_groups;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

// user register api response model. this removes password and status fields
// things that shoud not be returned in api response to the user
type UserResponse struct {
	ID        int    `json:"id" gorm:"primary_key"`
	FullName  string `json:"fullname"  binding:"required"`
	Email     string `json:"email" gorm:"unique"  binding:"required,email"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
