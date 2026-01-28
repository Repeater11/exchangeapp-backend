package app

import (
	"context"
	"errors"
	"exchangeapp/internal/repository"
	"time"
)

type LikeCountFlusher struct {
	counter  likeCounter
	writer   repository.ThreadLikeCountWriter
	batch    int
	interval time.Duration
}

type likeCounter interface {
	PopDirty(limit int) ([]uint, error)
	GetLikeCount(threadID uint) (int64, error)
	MarkDirty(threadID uint) error
}

func NewLikeCountFlusher(counter likeCounter, writer repository.ThreadLikeCountWriter, batch int, interval time.Duration) *LikeCountFlusher {
	if batch <= 0 {
		batch = 200
	}
	if interval <= 0 {
		interval = time.Second
	}

	return &LikeCountFlusher{
		counter:  counter,
		writer:   writer,
		batch:    batch,
		interval: interval,
	}
}

func (f *LikeCountFlusher) Run(ctx context.Context) {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.flushOnce()
		}
	}
}

func (f *LikeCountFlusher) flushOnce() {
	if f.batch <= 0 {
		return
	}

	ids, err := f.counter.PopDirty(f.batch)
	if err != nil {
		return
	}

	for _, id := range ids {
		val, err := f.counter.GetLikeCount(id)
		if err != nil {
			if !errors.Is(err, repository.ErrLikeCountNotFound) {
				_ = f.counter.MarkDirty(id)
			}
			continue
		}

		if err := f.writer.SetLikeCount(id, val); err != nil {
			_ = f.counter.MarkDirty(id)
		}
	}
}
