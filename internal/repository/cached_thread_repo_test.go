package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"exchangeapp/internal/models"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type fakeThreadRepoCache struct {
	findVal   *models.Thread
	findErr   error
	findCalls int
}

func (f *fakeThreadRepoCache) Create(*models.Thread) error {
	return nil
}

func (f *fakeThreadRepoCache) List(int, int) ([]models.Thread, error) {
	return nil, nil
}

func (f *fakeThreadRepoCache) ListAfter(time.Time, uint, int) ([]models.Thread, error) {
	return nil, nil
}

func (f *fakeThreadRepoCache) FindByID(id uint) (*models.Thread, error) {
	f.findCalls++
	return f.findVal, f.findErr
}

func (f *fakeThreadRepoCache) Count() (int64, error) {
	return 0, nil
}

func (f *fakeThreadRepoCache) ListByUserID(uint, int, int) ([]models.Thread, error) {
	return nil, nil
}

func (f *fakeThreadRepoCache) ListByUserIDAfter(uint, time.Time, uint, int) ([]models.Thread, error) {
	return nil, nil
}

func (f *fakeThreadRepoCache) CountByUserID(uint) (int64, error) {
	return 0, nil
}

func (f *fakeThreadRepoCache) Update(*models.Thread) error {
	return nil
}

func (f *fakeThreadRepoCache) DeleteByID(uint) error {
	return nil
}

func (f *fakeThreadRepoCache) IncrementLikeCount(uint, int) error {
	return nil
}

func (f *fakeThreadRepoCache) GetLikeCount(uint) (int64, error) {
	return 0, nil
}

type fakeRedisClient struct {
	val       []byte
	getErr    error
	setErr    error
	setKeys   []string
	setValues map[string]string
}

func (f *fakeRedisClient) Get(_ context.Context, _ string) *redis.StringCmd {
	cmd := redis.NewStringCmd(context.Background())
	cmd.SetVal(string(f.val))
	if f.getErr != nil {
		cmd.SetErr(f.getErr)
	}
	return cmd
}

func (f *fakeRedisClient) Set(_ context.Context, key string, value interface{}, _ time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(context.Background())
	if f.setErr != nil {
		cmd.SetErr(f.setErr)
	} else {
		f.setKeys = append(f.setKeys, key)
		if f.setValues == nil {
			f.setValues = make(map[string]string)
		}
		switch v := value.(type) {
		case []byte:
			f.setValues[key] = string(v)
		case string:
			f.setValues[key] = v
		}
	}
	return cmd
}

func (f *fakeRedisClient) Del(_ context.Context, _ ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(context.Background())
	cmd.SetVal(1)
	return cmd
}

func TestCachedThreadRepoFindByIDCacheHit(t *testing.T) {
	thread := &models.Thread{Model: gormModel(1), Title: "t1"}
	raw, _ := json.Marshal(thread)
	db := &fakeThreadRepoCache{}
	rdb := &fakeRedisClient{val: raw}
	repo := &CachedThreadRepo{
		db:  db,
		rdb: rdb,
		sf:  &singleflight.Group{},
	}

	got, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != 1 {
		t.Fatalf("expected cached thread id=1")
	}
	if db.findCalls != 0 {
		t.Fatalf("expected db not called, got %d", db.findCalls)
	}
}

func TestCachedThreadRepoFindByIDCacheMiss(t *testing.T) {
	thread := &models.Thread{Model: gormModel(2), Title: "t2"}
	db := &fakeThreadRepoCache{findVal: thread}
	rdb := &fakeRedisClient{getErr: redis.Nil}
	repo := &CachedThreadRepo{
		db:  db,
		rdb: rdb,
		sf:  &singleflight.Group{},
	}

	got, err := repo.FindByID(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != 2 {
		t.Fatalf("expected thread id=2")
	}
	if db.findCalls != 1 {
		t.Fatalf("expected db called once, got %d", db.findCalls)
	}
	if len(rdb.setKeys) != 1 {
		t.Fatalf("expected cache set once, got %d", len(rdb.setKeys))
	}
}

func TestCachedThreadRepoFindByIDNotFound(t *testing.T) {
	db := &fakeThreadRepoCache{findVal: nil}
	rdb := &fakeRedisClient{getErr: redis.Nil}
	repo := &CachedThreadRepo{
		db:  db,
		rdb: rdb,
		sf:  &singleflight.Group{},
	}

	got, err := repo.FindByID(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil thread")
	}
	if db.findCalls != 1 {
		t.Fatalf("expected db called once, got %d", db.findCalls)
	}
	if len(rdb.setKeys) != 1 {
		t.Fatalf("expected cache set once, got %d", len(rdb.setKeys))
	}
	if rdb.setValues[repo.cacheKey(3)] != threadCacheNotFound {
		t.Fatalf("expected not found cache set")
	}
}

func gormModel(id uint) gorm.Model {
	return gorm.Model{ID: id}
}
