package service

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

var ErrInvalidCredentials = errors.New("用户名或密码错误")

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

func (s *UserService) Login(req dto.LoginReq) (*dto.LoginResp, error) {
	u, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return &dto.LoginResp{
		Username: u.Username,
		ID:       u.ID,
	}, nil
}
