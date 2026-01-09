package service

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"testing"
)

type fakeReplyRepo struct {
	findResult *models.Reply
	findErr    error

	listResult  []models.Reply
	listErr     error
	countResult int64
	countErr    error

	updateErr error
	deleteErr error

	updated   *models.Reply
	deletedID uint
}

func (f *fakeReplyRepo) Create(*models.Reply) error {
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

func TestReplyServiceListByThreadID(t *testing.T) {
	svc := NewReplyService(
		&fakeReplyRepo{},
		&fakeThreadRepo{findResult: nil},
	)
	_, err := svc.ListByThreadID(1, 1, 10)
	if !errors.Is(err, ErrThreadNotFound) {
		t.Fatalf("expected ErrThreadNotFound, got %v", err)
	}

	replyRepo := &fakeReplyRepo{
		listResult:  []models.Reply{*reply(1, 2, 1)},
		countResult: 1,
	}
	svc = NewReplyService(replyRepo, &fakeThreadRepo{findResult: thread(1, 1)})
	resp, err := svc.ListByThreadID(1, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected result: %+v", resp)
	}
}

func TestReplyServiceUpdate(t *testing.T) {
	cases := []struct {
		name    string
		userID  uint
		reply   *models.Reply
		wantErr error
	}{
		{"not_found", 1, nil, ErrReplyNotFound},
		{"forbidden", 1, reply(1, 2, 1), ErrForbidden},
		{"ok", 1, reply(1, 1, 1), nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := &fakeReplyRepo{findResult: c.reply}
			svc := NewReplyService(repo, &fakeThreadRepo{findResult: thread(1, 1)})

			req := dto.UpdateReplyReq{Content: "new"}
			_, err := svc.Update(c.userID, 1, req)

			if !errors.Is(err, c.wantErr) {
				t.Fatalf("expected %v, got %v", c.wantErr, err)
			}
		})
	}
}

func TestReplyServiceDelete(t *testing.T) {
	repo := &fakeReplyRepo{findResult: reply(1, 1, 1)}
	svc := NewReplyService(repo, &fakeThreadRepo{findResult: thread(1, 1)})

	if err := svc.Delete(1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplyServiceListByUserID(t *testing.T) {
	repo := &fakeReplyRepo{
		listResult:  []models.Reply{*reply(1, 1, 1)},
		countResult: 1,
	}
	svc := NewReplyService(repo, &fakeThreadRepo{})

	resp, err := svc.ListByUserID(1, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected result: %+v", resp)
	}
}
