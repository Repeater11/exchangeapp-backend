package repository

import (
	"errors"
	"exchangeapp/internal/models"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrAlreadyLiked = errors.New("已点赞")
var ErrLikeNotFound = errors.New("未点赞")

type ThreadLikeRepository interface {
	Create(*models.ThreadLike) error
	Delete(userID, threadID uint) error
	Exists(userID, threadID uint) (bool, error)
}

type ThreadLikeRepo struct {
	db *gorm.DB
}

func NewThreadLikeRepository(db *gorm.DB) ThreadLikeRepository {
	return &ThreadLikeRepo{db: db}
}

func (r *ThreadLikeRepo) Create(t *models.ThreadLike) error {
	if err := r.db.Create(t).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyLiked
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return ErrAlreadyLiked
		}
		return fmt.Errorf("创建帖子点赞失败：%w", err)
	}
	return nil
}

func (r *ThreadLikeRepo) Delete(userID, threadID uint) error {
	res := r.db.Where("user_id = ? and thread_id = ?", userID, threadID).
		Delete(&models.ThreadLike{})
	if res.Error != nil {
		return fmt.Errorf("删除帖子点赞失败：%w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrLikeNotFound
	}
	return nil
}

func (r *ThreadLikeRepo) Exists(userID, threadID uint) (bool, error) {
	var cnt int64
	if err := r.db.Model(&models.ThreadLike{}).
		Where("user_id = ? and thread_id = ?", userID, threadID).
		Count(&cnt).Error; err != nil {
		return false, fmt.Errorf("查询帖子点赞失败：%w", err)
	}
	return cnt > 0, nil
}
