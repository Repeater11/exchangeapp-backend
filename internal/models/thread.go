package models

import (
	"time"

	"gorm.io/gorm"
)

type Thread struct {
	ID        uint      `gorm:"primaryKey;index:idx_threads_created_id,priority:2,sort:desc;index:idx_threads_user_created_id,priority:3,sort:desc"`
	CreatedAt time.Time `gorm:"index:idx_threads_created_id,priority:1,sort:desc;index:idx_threads_user_created_id,priority:2,sort:desc"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title     string
	Content   string
	UserID    uint  `gorm:"index:idx_threads_user_created_id,priority:1"`
	LikeCount int64 `gorm:"default:0"`
}
