package models

import "gorm.io/gorm"

type Thread struct {
	gorm.Model
	Title     string
	Content   string
	UserID    uint
	LikeCount int64 `gorm:"default:0"`
}
