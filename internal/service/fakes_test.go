package service

import (
	"exchangeapp/internal/models"

	"gorm.io/gorm"
)

type fakeThreadRepo struct {
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
	return nil, nil
}

func (f *fakeThreadRepo) Count() (int64, error) {
	return 0, nil
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

func gormModel(id uint) gorm.Model {
	return gorm.Model{
		ID: id,
	}
}

func thread(id, userID uint) *models.Thread {
	return &models.Thread{
		Model:  gormModel(id),
		UserID: userID,
	}
}

func reply(id, userID, threadID uint) *models.Reply {
	return &models.Reply{
		Model:    gormModel(id),
		UserID:   userID,
		ThreadID: threadID,
	}
}
