package models

import "time"

type GroupChat struct {
	ID        int     `gorm:"primary_key"`
	Title     string  `json:"title"  gorm:"unique"  binding:"required,uniq=group_chats.title"`
	Users     []*User `gorm:"many2many:user_groups;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type MessageStatus struct {
	Title  string `json:"title"`
	From   string `json:"from"  binding:"required"`
	Status string `json:"status"  binding:"required"`
}

type UserGroups struct {
	GroupChatID int
	UserID      int
}

type AddUser struct {
	GroupTitle string `json:"group_title"  binding:"required"`
	UserIDs    int    `json:"user_ids"  binding:"required"`
}
