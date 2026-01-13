package repository

import (
	"errors"
	"testing"
	"time"

	"exchangeapp/internal/models"
	"golang.org/x/sync/singleflight"
)

type fakeThreadRepo struct {
	getVal    int64
	getErr    error
	getCalls  int
	incErr    error
	incCalls  int
	lastDelta int
}

func (f *fakeThreadRepo) Create(*models.Thread) error {
	return nil
}

func (f *fakeThreadRepo) List(int, int) ([]models.Thread, error) {
	return nil, nil
}

func (f *fakeThreadRepo) FindByID(uint) (*models.Thread, error) {
	return nil, nil
}

func (f *fakeThreadRepo) Count() (int64, error) {
	return 0, nil
}

func (f *fakeThreadRepo) ListByUserID(uint, int, int) ([]models.Thread, error) {
	return nil, nil
}

func (f *fakeThreadRepo) CountByUserID(uint) (int64, error) {
	return 0, nil
}

func (f *fakeThreadRepo) Update(*models.Thread) error {
	return nil
}

func (f *fakeThreadRepo) DeleteByID(uint) error {
	return nil
}

func (f *fakeThreadRepo) IncrementLikeCount(threadID uint, delta int) error {
	f.incCalls++
	f.lastDelta = delta
	return f.incErr
}

func (f *fakeThreadRepo) GetLikeCount(threadID uint) (int64, error) {
	f.getCalls++
	return f.getVal, f.getErr
}

type fakeLikeCache struct {
	getVal    int64
	getErr    error
	setErr    error
	setCalls  int
	setVal    int64
	incErr    error
	incCalls  int
	lastDelta int
	lockErr   error
	lockCalls int
	unlockErr error
}

func (f *fakeLikeCache) IncrementLikeCount(threadID uint, delta int) error {
	f.incCalls++
	f.lastDelta = delta
	return f.incErr
}

func (f *fakeLikeCache) GetLikeCount(threadID uint) (int64, error) {
	return f.getVal, f.getErr
}

func (f *fakeLikeCache) setLikeCount(threadID uint, value int64) error {
	f.setCalls++
	f.setVal = value
	return f.setErr
}

func (f *fakeLikeCache) TryLockLikeCount(threadID uint, token string, ttl time.Duration) (bool, error) {
	f.lockCalls++
	if f.lockErr != nil {
		return false, f.lockErr
	}
	return true, nil
}

func (f *fakeLikeCache) UnlockLikeCount(threadID uint, token string) error {
	return f.unlockErr
}

func TestCachedThreadLikeCounterGetLikeCountCacheHit(t *testing.T) {
	db := &fakeThreadRepo{getVal: 7}
	cache := &fakeLikeCache{getVal: 9}
	counter := &CachedThreadLikeCounter{db: db, cache: cache, sf: &singleflight.Group{}}

	val, err := counter.GetLikeCount(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 9 {
		t.Fatalf("expected 9, got %d", val)
	}
	if db.getCalls != 0 {
		t.Fatalf("expected db not called, got %d", db.getCalls)
	}
}

func TestCachedThreadLikeCounterGetLikeCountCacheMiss(t *testing.T) {
	db := &fakeThreadRepo{getVal: 5}
	cache := &fakeLikeCache{getErr: ErrLikeCountNotFound}
	counter := &CachedThreadLikeCounter{db: db, cache: cache, sf: &singleflight.Group{}}

	val, err := counter.GetLikeCount(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 5 {
		t.Fatalf("expected 5, got %d", val)
	}
	if db.getCalls != 1 {
		t.Fatalf("expected db called once, got %d", db.getCalls)
	}
	if cache.setCalls != 1 || cache.setVal != 5 {
		t.Fatalf("expected cache set to 5, got calls=%d value=%d", cache.setCalls, cache.setVal)
	}
}

func TestCachedThreadLikeCounterGetLikeCountCacheErrorFallback(t *testing.T) {
	db := &fakeThreadRepo{getVal: 3}
	cache := &fakeLikeCache{getErr: errors.New("boom")}
	counter := &CachedThreadLikeCounter{db: db, cache: cache, sf: &singleflight.Group{}}

	val, err := counter.GetLikeCount(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 3 {
		t.Fatalf("expected 3, got %d", val)
	}
	if db.getCalls != 1 {
		t.Fatalf("expected db called once, got %d", db.getCalls)
	}
}

func TestCachedThreadLikeCounterIncrementLikeCountDBError(t *testing.T) {
	db := &fakeThreadRepo{incErr: errors.New("boom")}
	cache := &fakeLikeCache{}
	counter := &CachedThreadLikeCounter{db: db, cache: cache, sf: &singleflight.Group{}}

	if err := counter.IncrementLikeCount(1, 1); err == nil {
		t.Fatalf("expected error, got nil")
	}
	if cache.incCalls != 0 {
		t.Fatalf("expected cache not called, got %d", cache.incCalls)
	}
}

func TestCachedThreadLikeCounterIncrementLikeCountOK(t *testing.T) {
	db := &fakeThreadRepo{}
	cache := &fakeLikeCache{incErr: errors.New("cache")}
	counter := &CachedThreadLikeCounter{db: db, cache: cache, sf: &singleflight.Group{}}

	if err := counter.IncrementLikeCount(1, -1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db.incCalls != 1 || db.lastDelta != -1 {
		t.Fatalf("expected db delta -1, got calls=%d delta=%d", db.incCalls, db.lastDelta)
	}
	if cache.incCalls != 1 || cache.lastDelta != -1 {
		t.Fatalf("expected cache delta -1, got calls=%d delta=%d", cache.incCalls, cache.lastDelta)
	}
}
