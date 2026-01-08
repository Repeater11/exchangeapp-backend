package repository

import (
	"exchangeapp/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ReplyRepository interface {
	Create(*models.Reply) error
	ListByThreadID(threadID uint, limit, offset int) ([]models.Reply, error)
	CountByThreadID(threadID uint) (int64, error)
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

func (r *ReplyRepo) ListByThreadID(threadID uint, limit, offset int) ([]models.Reply, error) {
	var replies []models.Reply
	if err := r.db.Where("thread_id = ?", threadID).
		Order("created_at asc").
		Limit(limit).Offset(offset).
		Find(&replies).Error; err != nil {
		return nil, fmt.Errorf("查询回复失败：%w", err)
	}
	return replies, nil
}

func (r *ReplyRepo) CountByThreadID(threadID uint) (int64, error) {
	var total int64
	if err := r.db.Model(&models.Reply{}).
		Where("thread_id = ?", threadID).
		Count(&total).Error; err != nil {
		return 0, fmt.Errorf("统计回复失败：%w", err)
	}
	return total, nil
}
