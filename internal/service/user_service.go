package service

import (
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req dto.RegisterReq) (*dto.RegisterResp, error) {
	u := &models.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return &dto.RegisterResp{
		Username: u.Username,
		ID:       u.ID,
	}, nil
}
