package models

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content" binding:"required"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user"`
	Post Post `gorm:"foreignKey:PostID" json:"-"`
}