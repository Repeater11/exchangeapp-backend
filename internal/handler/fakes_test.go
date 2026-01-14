package handler

import (
	"time"

	"exchangeapp/internal/models"
)

type fakeThreadRepo struct {
	createErr    error
	listErr      error
	listAfterErr error
	countErr     error
	findErr      error
	updateErr    error
	deleteErr    error

	listResult         []models.Thread
	listAfterResult    []models.Thread
	listByUserAfterRes []models.Thread
	countResult        int64
	findResult         *models.Thread

	created   *models.Thread
	updated   *models.Thread
	deletedID uint
}

func (f *fakeThreadRepo) Create(t *models.Thread) error {
	if f.createErr != nil {
		return f.createErr
	}
	if t.ID == 0 {
		t.ID = 1
	}
	f.created = t
	return nil
}

func (f *fakeThreadRepo) List(limit, offset int) ([]models.Thread, error) {
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

func (f *fakeThreadRepo) FindByID(id uint) (*models.Thread, error) {
	return f.findResult, f.findErr
}

func (f *fakeThreadRepo) Update(t *models.Thread) error {
	f.updated = t
	return f.updateErr
}

func (f *fakeThreadRepo) DeleteByID(id uint) error {
	f.deletedID = id
	return f.deleteErr
}

func (f *fakeThreadRepo) IncrementLikeCount(threadID uint, delta int) error {
	return nil
}

func (f *fakeThreadRepo) GetLikeCount(threadID uint) (int64, error) {
	return 0, nil
}

type fakeReplyRepo struct {
	createErr error
	listErr   error
	countErr  error
	findErr   error
	updateErr error
	deleteErr error

	listResult  []models.Reply
	countResult int64
	findResult  *models.Reply

	created   *models.Reply
	updated   *models.Reply
	deletedID uint
}

func (f *fakeReplyRepo) Create(r *models.Reply) error {
	if f.createErr != nil {
		return f.createErr
	}
	if r.ID == 0 {
		r.ID = 1
	}
	f.created = r
	return nil
}

func (f *fakeReplyRepo) ListByThreadID(threadID uint, limit, offset int) ([]models.Reply, error) {
	return f.listResult, f.listErr
}

func (f *fakeReplyRepo) CountByThreadID(threadID uint) (int64, error) {
	return f.countResult, f.countErr
}

func (f *fakeReplyRepo) ListByUserID(userID uint, limit, offset int) ([]models.Reply, error) {
	return f.listResult, f.listErr
}

func (f *fakeReplyRepo) CountByUserID(userID uint) (int64, error) {
	return f.countResult, f.countErr
}

func (f *fakeReplyRepo) FindByID(id uint) (*models.Reply, error) {
	return f.findResult, f.findErr
}

func (f *fakeReplyRepo) Update(r *models.Reply) error {
	f.updated = r
	return f.updateErr
}

func (f *fakeReplyRepo) DeleteByID(id uint) error {
	f.deletedID = id
	return f.deleteErr
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
