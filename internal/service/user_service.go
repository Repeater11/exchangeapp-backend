package service

import (
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req dto.RegisterReq) (*dto.RegisterResp, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("加密时出错：%w", err)
	}
	u := &models.User{
		Username: req.Username,
		Password: string(hashed),
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return &dto.RegisterResp{
		Username: u.Username,
		ID:       u.ID,
	}, nil
}
