package service

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"testing"
	"time"

	"gorm.io/gorm"
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
			svc := NewThreadService(repo, &fakeThreadLikeRepo{}, repo)

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
	svc := NewThreadService(repo, &fakeThreadLikeRepo{}, repo)

	if err := svc.Delete(1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestThreadServiceListByUserID(t *testing.T) {
	repo := &fakeThreadRepo{
		listResult:  []models.Thread{*thread(1, 1)},
		countResult: 1,
	}
	svc := NewThreadService(repo, &fakeThreadLikeRepo{}, repo)

	resp, err := svc.ListByUserID(1, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected result: %+v", resp)
	}
}

func TestThreadServiceListAfter(t *testing.T) {
	ts := time.Unix(0, 123)
	repo := &fakeThreadRepo{
		listAfterResult: []models.Thread{
			{Model: gorm.Model{ID: 7, CreatedAt: ts}, Title: "t1", UserID: 1},
		},
	}
	svc := NewThreadService(repo, &fakeThreadLikeRepo{}, repo)

	resp, err := svc.ListAfter(time.Unix(0, 1), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("unexpected result: %+v", resp)
	}
	wantCursor := "123_7"
	if resp.NextCursor != wantCursor {
		t.Fatalf("expected next_cursor %s, got %s", wantCursor, resp.NextCursor)
	}
}

func TestThreadServiceListByUserIDAfter(t *testing.T) {
	ts := time.Unix(0, 456)
	repo := &fakeThreadRepo{
		listByUserAfterRes: []models.Thread{
			{Model: gorm.Model{ID: 9, CreatedAt: ts}, Title: "t2", UserID: 2},
		},
	}
	svc := NewThreadService(repo, &fakeThreadLikeRepo{}, repo)

	resp, err := svc.ListByUserIDAfter(2, time.Unix(0, 1), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("unexpected result: %+v", resp)
	}
	wantCursor := "456_9"
	if resp.NextCursor != wantCursor {
		t.Fatalf("expected next_cursor %s, got %s", wantCursor, resp.NextCursor)
	}
}
