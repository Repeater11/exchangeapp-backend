package models

import "gorm.io/gorm"

type Reply struct {
	gorm.Model
	ThreadID uint
	Content  string
	UserID   uint
}
