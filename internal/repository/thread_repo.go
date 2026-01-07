package repository

import (
	"errors"
	"exchangeapp/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ThreadRepository interface {
	Create(*models.Thread) error
	List() ([]models.Thread, error)
	FindByID(id uint) (*models.Thread, error)
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

func (r *ThreadRepo) List() ([]models.Thread, error) {
	var threads []models.Thread
	if err := r.db.Order("created_at desc").Find(&threads).Error; err != nil {
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
