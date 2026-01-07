package service

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
)

var ErrThreadNotFound = errors.New("帖子不存在")

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

func (s *ReplyService) ListByThreadID(threadID uint) ([]dto.ReplyResp, error) {
	rs, err := s.replyRepo.ListByThreadID(threadID)
	if err != nil {
		return nil, err
	}

	rrs := make([]dto.ReplyResp, len(rs))
	for i := range rrs {
		rrs[i] = dto.ReplyResp{
			ID:        rs[i].ID,
			ThreadID:  rs[i].ThreadID,
			Content:   rs[i].Content,
			UserID:    rs[i].UserID,
			CreatedAt: rs[i].CreatedAt,
		}
	}
	return rrs, nil
}
