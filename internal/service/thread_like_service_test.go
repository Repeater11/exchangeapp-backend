package service

import (
	"errors"
	"exchangeapp/internal/repository"
	"testing"
)

func TestThreadLikeServiceLike(t *testing.T) {
	cases := []struct {
		name    string
		thread  bool
		repoErr error
		wantErr error
	}{
		{"not_found", false, nil, ErrThreadNotFound},
		{"already_liked", true, repository.ErrAlreadyLiked, repository.ErrAlreadyLiked},
		{"ok", true, nil, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			threadRepo := &fakeThreadRepo{}
			if c.thread {
				threadRepo.findResult = thread(1, 1)
			}
			likeRepo := &fakeThreadLikeRepo{createErr: c.repoErr}

			svc := NewThreadLikeService(threadRepo, likeRepo)
			err := svc.Like(1, 1)

			if !errors.Is(err, c.wantErr) {
				t.Fatalf("expected %v, got %v", c.wantErr, err)
			}
		})
	}
}

func TestThreadLikeServiceUnlike(t *testing.T) {
	cases := []struct {
		name    string
		thread  bool
		repoErr error
		wantErr error
	}{
		{"not_found", false, nil, ErrThreadNotFound},
		{"not_liked", true, repository.ErrLikeNotFound, repository.ErrLikeNotFound},
		{"ok", true, nil, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			threadRepo := &fakeThreadRepo{}
			if c.thread {
				threadRepo.findResult = thread(1, 1)
			}
			likeRepo := &fakeThreadLikeRepo{deleteErr: c.repoErr}

			svc := NewThreadLikeService(threadRepo, likeRepo)
			err := svc.Unlike(1, 1)

			if !errors.Is(err, c.wantErr) {
				t.Fatalf("expected %v, got %v", c.wantErr, err)
			}
		})
	}
}

func TestThreadLikeServiceIsLiked(t *testing.T) {
	errBoom := errors.New("boom")
	cases := []struct {
		name     string
		thread   bool
		exists   bool
		repoErr  error
		wantErr  error
		wantLike bool
	}{
		{"not_found", false, false, nil, ErrThreadNotFound, false},
		{"repo_error", true, false, errBoom, errBoom, false},
		{"liked", true, true, nil, nil, true},
		{"not_liked", true, false, nil, nil, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			threadRepo := &fakeThreadRepo{}
			if c.thread {
				threadRepo.findResult = thread(1, 1)
			}
			likeRepo := &fakeThreadLikeRepo{
				exists:    c.exists,
				existsErr: c.repoErr,
			}

			svc := NewThreadLikeService(threadRepo, likeRepo)
			got, err := svc.IsLiked(1, 1)

			if !errors.Is(err, c.wantErr) {
				t.Fatalf("expected %v, got %v", c.wantErr, err)
			}
			if err == nil && got != c.wantLike {
				t.Fatalf("expected %v, got %v", c.wantLike, got)
			}
		})
	}
}
