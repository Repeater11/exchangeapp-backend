package repository

import (
	"errors"
	"exchangeapp/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ReplyRepository interface {
	Create(*models.Reply) error
	ListByThreadID(threadID uint, limit, offset int) ([]models.Reply, error)
	CountByThreadID(threadID uint) (int64, error)
	ListByUserID(userID uint, limit, offset int) ([]models.Reply, error)
	CountByUserID(userID uint) (int64, error)
	FindByID(id uint) (*models.Reply, error)
	Update(*models.Reply) error
	DeleteByID(id uint) error
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

func (r *ReplyRepo) ListByUserID(userID uint, limit, offset int) ([]models.Reply, error) {
	var replies []models.Reply
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).Offset(offset).
		Find(&replies).Error; err != nil {
		return nil, fmt.Errorf("查询回复失败：%w", err)
	}
	return replies, nil
}

func (r *ReplyRepo) CountByUserID(userID uint) (int64, error) {
	var total int64
	if err := r.db.Model(&models.Reply{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return 0, fmt.Errorf("统计回复失败：%w", err)
	}
	return total, nil
}

func (r *ReplyRepo) FindByID(id uint) (*models.Reply, error) {
	var rp models.Reply
	if err := r.db.First(&rp, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询评论失败：%w", err)
	}
	return &rp, nil
}

func (r *ReplyRepo) Update(rp *models.Reply) error {
	if err := r.db.Model(&models.Reply{}).
		Where("id = ?", rp.ID).
		Updates(map[string]interface{}{
			"content": rp.Content,
		}).Error; err != nil {
		return fmt.Errorf("更新评论失败：%w", err)
	}
	return nil
}

func (r *ReplyRepo) DeleteByID(id uint) error {
	if err := r.db.Delete(&models.Reply{}, id).Error; err != nil {
		return fmt.Errorf("删除评论失败：%w", err)
	}
	return nil
}
