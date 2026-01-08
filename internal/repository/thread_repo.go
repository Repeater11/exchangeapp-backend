package repository

import (
	"errors"
	"exchangeapp/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ThreadRepository interface {
	Create(*models.Thread) error
	List(limit, offset int) ([]models.Thread, error)
	FindByID(id uint) (*models.Thread, error)
	Count() (int64, error)
	Update(*models.Thread) error
	DeleteByID(id uint) error
}

type ThreadRepo struct {
	db *gorm.DB
}

func NewThreadRepository(db *gorm.DB) ThreadRepository {
	return &ThreadRepo{db: db}
}

func (r *ThreadRepo) Create(t *models.Thread) error {
	if err := r.db.Create(t).Error; err != nil {
		return fmt.Errorf("创建帖子失败：%w", err)
	}
	return nil
}

func (r *ThreadRepo) List(limit, offset int) ([]models.Thread, error) {
	var threads []models.Thread
	if err := r.db.Order("created_at desc").
		Limit(limit).Offset(offset).
		Find(&threads).Error; err != nil {
		return nil, fmt.Errorf("查询帖子失败：%w", err)
	}
	return threads, nil
}

func (r *ThreadRepo) FindByID(id uint) (*models.Thread, error) {
	var t models.Thread
	if err := r.db.First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询帖子失败：%w", err)
	}
	return &t, nil
}

func (r *ThreadRepo) Count() (int64, error) {
	var total int64
	if err := r.db.Model(&models.Thread{}).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("统计帖子失败：%w", err)
	}

	return total, nil
}

func (r *ThreadRepo) Update(t *models.Thread) error {
	if err := r.db.Model(&models.Thread{}).
		Where("id = ?", t.ID).
		Updates(map[string]interface{}{
			"title":   t.Title,
			"content": t.Content,
		}).Error; err != nil {
		return fmt.Errorf("更新帖子失败：%w", err)
	}
	return nil
}

func (r *ThreadRepo) DeleteByID(id uint) error {
	if err := r.db.Delete(&models.Thread{}, id).Error; err != nil {
		return fmt.Errorf("删除帖子失败：%w", err)
	}
	return nil
}
