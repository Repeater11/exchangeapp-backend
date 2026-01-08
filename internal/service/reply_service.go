package service

import (
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
)

type ReplyService struct {
	replyRepo  repository.ReplyRepository
	threadRepo repository.ThreadRepository
}

func NewReplyService(replyRepo repository.ReplyRepository,
	threadRepo repository.ThreadRepository) *ReplyService {
	return &ReplyService{
		replyRepo:  replyRepo,
		threadRepo: threadRepo,
	}
}

func (s *ReplyService) Create(userID uint, threadID uint, req dto.CreateReplyReq) (*dto.ReplyResp, error) {
	t, err := s.threadRepo.FindByID(threadID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrThreadNotFound
	}

	r := &models.Reply{
		ThreadID: threadID,
		Content:  req.Content,
		UserID:   userID,
	}

	if err := s.replyRepo.Create(r); err != nil {
		return nil, err
	}

	return &dto.ReplyResp{
		ID:        r.ID,
		ThreadID:  r.ThreadID,
		Content:   r.Content,
		UserID:    r.UserID,
		CreatedAt: r.CreatedAt,
	}, nil
}

func (s *ReplyService) ListByThreadID(threadID uint, page, size int) (*dto.ReplyListResp, error) {
	t, err := s.threadRepo.FindByID(threadID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrThreadNotFound
	}

	offset := (page - 1) * size

	total, err := s.replyRepo.CountByThreadID(threadID)
	if err != nil {
		return nil, err
	}
	rs, err := s.replyRepo.ListByThreadID(threadID, size, offset)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReplyResp, len(rs))
	for i := range rs {
		items[i] = dto.ReplyResp{
			ID:        rs[i].ID,
			ThreadID:  rs[i].ThreadID,
			Content:   rs[i].Content,
			UserID:    rs[i].UserID,
			CreatedAt: rs[i].CreatedAt,
		}
	}

	return &dto.ReplyListResp{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

func (s *ReplyService) Update(userID, id uint, req dto.UpdateReplyReq) (*dto.ReplyResp, error) {
	r, err := s.replyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, ErrReplyNotFound
	}
	if r.UserID != userID {
		return nil, ErrForbidden
	}

	r.Content = req.Content

	if err := s.replyRepo.Update(r); err != nil {
		return nil, err
	}
	return &dto.ReplyResp{
		ID:        r.ID,
		ThreadID:  r.ThreadID,
		Content:   r.Content,
		UserID:    r.UserID,
		CreatedAt: r.CreatedAt,
	}, nil
}

func (s *ReplyService) Delete(userID, id uint) error {
	r, err := s.replyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if r == nil {
		return ErrReplyNotFound
	}
	if r.UserID != userID {
		return ErrForbidden
	}

	return s.replyRepo.DeleteByID(id)
}
