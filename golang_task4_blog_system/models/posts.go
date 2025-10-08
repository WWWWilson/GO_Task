package models

import (
	"time"
)

type Post struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title" binding:"required"`
	Content   string    `gorm:"type:text;not null" json:"content" binding:"required"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;" json:"comments,omitempty"`
}
