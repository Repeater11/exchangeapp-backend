package repository

import (
	"errors"
	"exchangeapp/internal/models"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(*models.User) error
	FindByUsername(username string) (*models.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
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

func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败：%w", err)
	}
	return &u, nil
}
