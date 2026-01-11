package service

import (
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"

	"gorm.io/gorm"
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
	txer, ok1 := s.threadRepo.(repository.Transactioner)
	trWithTx, ok2 := s.threadRepo.(repository.ThreadRepoWithTx)
	lrWithTx, ok3 := s.likeRepo.(repository.ThreadLikeRepoWithTx)

	if ok1 && ok2 && ok3 {
		return txer.Transaction(func(tx *gorm.DB) error {
			tr := trWithTx.WithTx(tx)
			lr := lrWithTx.WithTx(tx)

			t, err := tr.FindByID(threadID)
			if err != nil {
				return err
			}
			if t == nil {
				return ErrThreadNotFound
			}

			if err := lr.Create(&models.ThreadLike{
				UserID:   userID,
				ThreadID: threadID,
			}); err != nil {
				return err
			}
			return tr.IncrementLikeCount(threadID, 1)
		})
	}

	t, err := s.threadRepo.FindByID(threadID)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrThreadNotFound
	}

	if err := s.likeRepo.Create(&models.ThreadLike{
		UserID:   userID,
		ThreadID: threadID,
	}); err != nil {
		return err
	}

	return s.threadRepo.IncrementLikeCount(threadID, 1)
}

func (s *ThreadLikeService) Unlike(userID, threadID uint) error {
	txer, ok1 := s.threadRepo.(repository.Transactioner)
	trWithTx, ok2 := s.threadRepo.(repository.ThreadRepoWithTx)
	lrWithTx, ok3 := s.likeRepo.(repository.ThreadLikeRepoWithTx)

	if ok1 && ok2 && ok3 {
		return txer.Transaction(func(tx *gorm.DB) error {
			tr := trWithTx.WithTx(tx)
			lr := lrWithTx.WithTx(tx)

			t, err := tr.FindByID(threadID)
			if err != nil {
				return err
			}
			if t == nil {
				return ErrThreadNotFound
			}

			if err := lr.Delete(userID, threadID); err != nil {
				return err
			}
			return tr.IncrementLikeCount(threadID, -1)
		})
	}

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
	return s.threadRepo.IncrementLikeCount(threadID, -1)
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
