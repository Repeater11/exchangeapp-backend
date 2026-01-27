package models

import (
	"time"

	"gorm.io/gorm"
)

type Reply struct {
	ID        uint      `gorm:"primaryKey;index:idx_replies_thread_created_id,priority:3,sort:desc;index:idx_replies_user_created_id,priority:3,sort:desc"`
	CreatedAt time.Time `gorm:"index:idx_replies_thread_created_id,priority:2,sort:desc;index:idx_replies_user_created_id,priority:2,sort:desc"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ThreadID uint `gorm:"index:idx_replies_thread_created_id,priority:1"`
	Content  string
	UserID   uint `gorm:"index:idx_replies_user_created_id,priority:1"`
}
