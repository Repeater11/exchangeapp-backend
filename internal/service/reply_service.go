package service

import (
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
	"fmt"
	"time"
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

	next := ""
	if len(rs) > 0 {
		last := rs[len(rs)-1]
		next = fmt.Sprintf("%d_%d", last.CreatedAt.UnixNano(), last.ID)
	}

	return &dto.ReplyListResp{
		Items:      items,
		Total:      total,
		Page:       page,
		Size:       size,
		NextCursor: next,
	}, nil
}

func (s *ReplyService) ListByThreadIDAfter(threadID uint, cursorTime time.Time, cursorID uint, size int) (*dto.ReplyListResp, error) {
	t, err := s.threadRepo.FindByID(threadID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrThreadNotFound
	}

	replies, err := s.replyRepo.ListByThreadIDAfter(threadID, cursorTime, cursorID, size)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReplyResp, len(replies))
	for i := range replies {
		items[i] = dto.ReplyResp{
			ID:        replies[i].ID,
			ThreadID:  replies[i].ThreadID,
			Content:   replies[i].Content,
			UserID:    replies[i].UserID,
			CreatedAt: replies[i].CreatedAt,
		}
	}

	next := ""
	if len(replies) > 0 {
		last := replies[len(replies)-1]
		next = fmt.Sprintf("%d_%d", last.CreatedAt.UnixNano(), last.ID)
	}

	return &dto.ReplyListResp{
		Items:      items,
		Total:      0,
		Page:       0,
		Size:       size,
		NextCursor: next,
	}, nil
}

func (s *ReplyService) ListByUserID(userID uint, page, size int) (*dto.ReplyListResp, error) {
	offset := (page - 1) * size

	total, err := s.replyRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}
	rs, err := s.replyRepo.ListByUserID(userID, size, offset)
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

	next := ""
	if len(rs) > 0 {
		last := rs[len(rs)-1]
		next = fmt.Sprintf("%d_%d", last.CreatedAt.UnixNano(), last.ID)
	}

	return &dto.ReplyListResp{
		Items:      items,
		Total:      total,
		Page:       page,
		Size:       size,
		NextCursor: next,
	}, nil
}

func (s *ReplyService) ListByUserIDAfter(userID uint, cursorTime time.Time, cursorID uint, size int) (*dto.ReplyListResp, error) {
	replies, err := s.replyRepo.ListByUserIDAfter(userID, cursorTime, cursorID, size)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReplyResp, len(replies))
	for i := range replies {
		items[i] = dto.ReplyResp{
			ID:        replies[i].ID,
			ThreadID:  replies[i].ThreadID,
			Content:   replies[i].Content,
			UserID:    replies[i].UserID,
			CreatedAt: replies[i].CreatedAt,
		}
	}

	next := ""
	if len(replies) > 0 {
		last := replies[len(replies)-1]
		next = fmt.Sprintf("%d_%d", last.CreatedAt.UnixNano(), last.ID)
	}

	return &dto.ReplyListResp{
		Items:      items,
		Total:      0,
		Page:       0,
		Size:       size,
		NextCursor: next,
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
