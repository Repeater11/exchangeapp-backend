package models

import "gorm.io/gorm"

type ThreadLike struct {
	gorm.Model
	UserID   uint `gorm:"uniqueIndex:uidx_user_thread"`
	ThreadID uint `gorm:"uniqueIndex:uidx_user_thread"`
}
