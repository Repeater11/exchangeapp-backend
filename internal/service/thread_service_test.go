package service

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"testing"
)

func TestThreadServiceUpdate(t *testing.T) {
	cases := []struct {
		name    string
		userID  uint
		thread  *models.Thread
		wantErr error
	}{
		{"not_found", 1, nil, ErrThreadNotFound},
		{"forbidden", 1, thread(1, 2), ErrForbidden},
		{"ok", 1, thread(1, 1), nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := &fakeThreadRepo{findResult: c.thread}
			svc := NewThreadService(repo, &fakeThreadLikeRepo{})

			req := dto.UpdateThreadReq{Title: "t", Content: "c"}
			_, err := svc.Update(c.userID, 1, req)

			if !errors.Is(err, c.wantErr) {
				t.Fatalf("expected %v, got %v", c.wantErr, err)
			}
		})
	}
}

func TestThreadServiceDelete(t *testing.T) {
	repo := &fakeThreadRepo{
		findResult: thread(1, 1),
	}
	svc := NewThreadService(repo, &fakeThreadLikeRepo{})

	if err := svc.Delete(1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestThreadServiceListByUserID(t *testing.T) {
	repo := &fakeThreadRepo{
		listResult:  []models.Thread{*thread(1, 1)},
		countResult: 1,
	}
	svc := NewThreadService(repo, &fakeThreadLikeRepo{})

	resp, err := svc.ListByUserID(1, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected result: %+v", resp)
	}
}
