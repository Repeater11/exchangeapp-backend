package service

import (
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
)

type ThreadService struct {
	repo     repository.ThreadRepository
	likeRepo repository.ThreadLikeRepository
}

func NewThreadService(repo repository.ThreadRepository, likeRepo repository.ThreadLikeRepository) *ThreadService {
	return &ThreadService{repo: repo, likeRepo: likeRepo}
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

func (s *ThreadService) List(page, size int) (*dto.ThreadListResp, error) {
	offset := (page - 1) * size

	total, err := s.repo.Count()
	if err != nil {
		return nil, err
	}
	ts, err := s.repo.List(size, offset)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ThreadSummaryResp, len(ts))
	for i := range ts {
		items[i] = dto.ThreadSummaryResp{
			ID:        ts[i].ID,
			Title:     ts[i].Title,
			UserID:    ts[i].UserID,
			CreatedAt: ts[i].CreatedAt,
		}
	}

	return &dto.ThreadListResp{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

func (s *ThreadService) ListByUserID(userID uint, page, size int) (*dto.ThreadListResp, error) {
	offset := (page - 1) * size

	total, err := s.repo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}
	ts, err := s.repo.ListByUserID(userID, size, offset)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ThreadSummaryResp, len(ts))
	for i := range ts {
		items[i] = dto.ThreadSummaryResp{
			ID:        ts[i].ID,
			Title:     ts[i].Title,
			UserID:    ts[i].UserID,
			CreatedAt: ts[i].CreatedAt,
		}
	}

	return &dto.ThreadListResp{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

func (s *ThreadService) GetByID(id uint) (*dto.ThreadDetailResp, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrThreadNotFound
	}

	likeCount, err := s.likeRepo.CountByThreadID(id)
	if err != nil {
		return nil, err
	}

	return &dto.ThreadDetailResp{
		ID:        t.ID,
		Title:     t.Title,
		Content:   t.Content,
		UserID:    t.UserID,
		LikeCount: likeCount,
		CreatedAt: t.CreatedAt,
	}, nil
}

func (s *ThreadService) Update(userID, id uint, req dto.UpdateThreadReq) (*dto.ThreadDetailResp, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrThreadNotFound
	}
	if t.UserID != userID {
		return nil, ErrForbidden
	}

	t.Title = req.Title
	t.Content = req.Content

	if err := s.repo.Update(t); err != nil {
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

func (s *ThreadService) Delete(userID, id uint) error {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrThreadNotFound
	}
	if t.UserID != userID {
		return ErrForbidden
	}

	return s.repo.DeleteByID(id)
}
