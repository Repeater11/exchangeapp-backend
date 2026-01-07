package service

import (
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
)

type ThreadService struct {
	repo repository.ThreadRepository
}

func NewThreadService(repo repository.ThreadRepository) *ThreadService {
	return &ThreadService{repo: repo}
}

func (s *ThreadService) Create(userID uint, req dto.CreateThreadReq) (*dto.ThreadDetailResp, error) {
	t := &models.Thread{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.repo.Create(t); err != nil {
		return nil, err
	}

	return &dto.ThreadDetailResp{
		ID:        t.ID,
		Title:     t.Title,
		Content:   t.Content,
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
	}, nil
}

func (s *ThreadService) List() ([]dto.ThreadSummaryResp, error) {
	ts, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	trs := make([]dto.ThreadSummaryResp, len(ts))
	for i := range ts {
		trs[i] = dto.ThreadSummaryResp{
			ID:        ts[i].ID,
			Title:     ts[i].Title,
			UserID:    ts[i].UserID,
			CreatedAt: ts[i].CreatedAt,
		}
	}
	return trs, nil
}
