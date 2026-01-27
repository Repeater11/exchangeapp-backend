package service

import (
	"time"

	"exchangeapp/internal/models"

	"gorm.io/gorm"
)

type fakeThreadRepo struct {
	listResult         []models.Thread
	listErr            error
	listAfterResult    []models.Thread
	listAfterErr       error
	listByUserAfterRes []models.Thread
	countResult        int64
	countErr           error

	findResult *models.Thread
	findErr    error

	updateErr error
	deleteErr error

	updated  *models.Thread
	deleteID uint
}

func (f *fakeThreadRepo) Create(*models.Thread) error {
	return nil
}

func (f *fakeThreadRepo) List(int, int) ([]models.Thread, error) {
	return f.listResult, f.listErr
}

func (f *fakeThreadRepo) ListAfter(cursorTime time.Time, cursorID uint, limit int) ([]models.Thread, error) {
	if f.listAfterResult != nil || f.listAfterErr != nil {
		return f.listAfterResult, f.listAfterErr
	}
	return f.listResult, f.listErr
}

func (f *fakeThreadRepo) Count() (int64, error) {
	return f.countResult, f.countErr
}

func (f *fakeThreadRepo) ListByUserID(userID uint, limit, offset int) ([]models.Thread, error) {
	return f.listResult, f.listErr
}

func (f *fakeThreadRepo) ListByUserIDAfter(userID uint, cursorTime time.Time, cursorID uint, limit int) ([]models.Thread, error) {
	if f.listByUserAfterRes != nil || f.listAfterErr != nil {
		return f.listByUserAfterRes, f.listAfterErr
	}
	return f.listResult, f.listErr
}

func (f *fakeThreadRepo) CountByUserID(userID uint) (int64, error) {
	return f.countResult, f.countErr
}

func (f *fakeThreadRepo) FindByID(uint) (*models.Thread, error) {
	return f.findResult, f.findErr
}

func (f *fakeThreadRepo) Update(t *models.Thread) error {
	f.updated = t
	return f.updateErr
}

func (f *fakeThreadRepo) DeleteByID(id uint) error {
	f.deleteID = id
	return f.deleteErr
}

func (f *fakeThreadRepo) IncrementLikeCount(threadID uint, delta int) error {
	return nil
}

func (f *fakeThreadRepo) GetLikeCount(threadID uint) (int64, error) {
	return 0, nil
}

type fakeThreadLikeRepo struct {
	createErr error
	deleteErr error
	existsErr error
	exists    bool

	created *models.ThreadLike
	deleted bool
}

func (f *fakeThreadLikeRepo) Create(t *models.ThreadLike) error {
	f.created = t
	return f.createErr
}

func (f *fakeThreadLikeRepo) Delete(userID, threadID uint) error {
	f.deleted = true
	return f.deleteErr
}

func (f *fakeThreadLikeRepo) Exists(userID, threadID uint) (bool, error) {
	return f.exists, f.existsErr
}

func (f *fakeThreadLikeRepo) CountByThreadID(threadID uint) (int64, error) {
	return 0, nil
}

func gormModel(id uint) gorm.Model {
	return gorm.Model{
		ID: id,
	}
}

func thread(id, userID uint) *models.Thread {
	return &models.Thread{
		ID:     id,
		UserID: userID,
	}
}

func reply(id, userID, threadID uint) *models.Reply {
	return &models.Reply{
		ID:       id,
		UserID:   userID,
		ThreadID: threadID,
	}
}
