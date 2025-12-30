package repository

import (
	"errors"
	"exchangeapp/internal/models"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

var ErrUserExists = errors.New("用户名已经存在")

func (r *UserRepo) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrUserExists
		}

		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return ErrUserExists
		}

		return fmt.Errorf("创建用户失败：%w", err)
	}
	return nil
}
