package service

import (
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
)

type ThreadLikeService struct {
	threadRepo repository.ThreadRepository
	likeRepo   repository.ThreadLikeRepository
}

func NewThreadLikeService(threadRepo repository.ThreadRepository, likeRepo repository.ThreadLikeRepository) *ThreadLikeService {
	return &ThreadLikeService{
		threadRepo: threadRepo,
		likeRepo:   likeRepo,
	}
}

func (s *ThreadLikeService) Like(userID, threadID uint) error {
	t, err := s.threadRepo.FindByID(threadID)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrThreadNotFound
	}

	tl := &models.ThreadLike{
		UserID:   userID,
		ThreadID: threadID,
	}
	if err := s.likeRepo.Create(tl); err != nil {
		return err
	}
	if err := s.threadRepo.IncrementLikeCount(threadID, 1); err != nil {
		return err
	}
	return nil
}

func (s *ThreadLikeService) Unlike(userID, threadID uint) error {
	t, err := s.threadRepo.FindByID(threadID)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrThreadNotFound
	}

	if err := s.likeRepo.Delete(userID, threadID); err != nil {
		return err
	}
	if err := s.threadRepo.IncrementLikeCount(threadID, -1); err != nil {
		return err
	}
	return nil
}

func (s *ThreadLikeService) IsLiked(userID, threadID uint) (bool, error) {
	t, err := s.threadRepo.FindByID(threadID)
	if err != nil {
		return false, err
	}
	if t == nil {
		return false, ErrThreadNotFound
	}
	return s.likeRepo.Exists(userID, threadID)
}
