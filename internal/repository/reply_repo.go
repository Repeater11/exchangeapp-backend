package repository

import (
	"exchangeapp/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ReplyRepository interface {
	Create(*models.Reply) error
	ListByThreadID(threadID uint) ([]models.Reply, error)
}

type ReplyRepo struct {
	db *gorm.DB
}

func NewReplyRepository(db *gorm.DB) ReplyRepository {
	return &ReplyRepo{db: db}
}

func (r *ReplyRepo) Create(reply *models.Reply) error {
	if err := r.db.Create(reply).Error; err != nil {
		return fmt.Errorf("创建回复失败：%w", err)
	}
	return nil
}

func (r *ReplyRepo) ListByThreadID(threadID uint) ([]models.Reply, error) {
	var replies []models.Reply
	if err := r.db.Where("thread_id = ?", threadID).
		Order("created_at asc").
		Find(&replies).Error; err != nil {
		return nil, fmt.Errorf("查询回复失败：%w", err)
	}
	return replies, nil
}
