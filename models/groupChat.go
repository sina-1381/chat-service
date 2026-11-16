package models

import "time"

// group chat db model . this also is used for create group chat api validation
type GroupChat struct {
	ID        int     `gorm:"primary_key"`
	Title     string  `json:"title"  gorm:"unique"  binding:"required"`
	Users     []*User `gorm:"many2many:user_groups;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// message status update api validation
type MessageStatus struct {
	Title  string `json:"title"`
	From   string `json:"from"  binding:"required"`
	Status string `json:"status"  binding:"required"`
}

// user group chat many2many table model
type UserGroups struct {
	GroupChatID int
	UserID      int
}

// add user to group api validation
type AddUser struct {
	GroupTitle string `json:"group_title"  binding:"required"`
	UserIDs    int    `json:"user_ids"  binding:"required"`
}
