package repository

import (
	"context"
	"encoding/json"
	"exchangeapp/internal/models"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	threadCacheTTL       = 10 * time.Minute
	threadCacheTTLJitter = 1 * time.Minute

	threadCacheNotFound    = "__not_found__"
	threadCacheNotFoundTTL = 30 * time.Second
)

type threadCacheStore interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type CachedThreadRepo struct {
	db  ThreadRepository
	rdb threadCacheStore
	sf  *singleflight.Group
}

func NewCachedThreadRepository(db ThreadRepository, rdb *redis.Client) *CachedThreadRepo {
	return &CachedThreadRepo{
		db:  db,
		rdb: rdb,
		sf:  &singleflight.Group{},
	}
}

func (c *CachedThreadRepo) cacheKey(id uint) string {
	return fmt.Sprintf("thread:detail:%d", id)
}

func (c *CachedThreadRepo) getCache(id uint) (*models.Thread, bool, error) {
	val, err := c.rdb.Get(context.Background(), c.cacheKey(id)).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	if string(val) == threadCacheNotFound {
		return nil, true, nil
	}

	var t models.Thread
	if err := json.Unmarshal(val, &t); err != nil {
		return nil, false, err
	}
	return &t, true, nil
}

func (c *CachedThreadRepo) setCache(t *models.Thread) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	jitter := time.Duration(time.Now().UnixNano() % int64(threadCacheTTLJitter))
	ttl := threadCacheTTL + jitter
	return c.rdb.Set(context.Background(), c.cacheKey(t.ID), b, ttl).Err()
}

func (c *CachedThreadRepo) deleteCache(id uint) {
	_ = c.rdb.Del(context.Background(), c.cacheKey(id)).Err()
}

func (c *CachedThreadRepo) FindByID(id uint) (*models.Thread, error) {
	if t, hit, err := c.getCache(id); err == nil && hit {
		return t, nil
	}

	key := c.cacheKey(id)
	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		t, err := c.db.FindByID(id)
		if err != nil || t == nil {
			return t, err
		}
		_ = c.setCache(t)
		return t, nil
	})
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	t, _ := v.(*models.Thread)
	if t == nil {
		_ = c.rdb.Set(context.Background(), c.cacheKey(id), threadCacheNotFound, threadCacheNotFoundTTL).Err()
		return nil, nil
	}
	return t, nil
}

func (c *CachedThreadRepo) Create(t *models.Thread) error {
	if err := c.db.Create(t); err != nil {
		return err
	}
	_ = c.setCache(t)
	return nil
}

func (c *CachedThreadRepo) Update(t *models.Thread) error {
	if err := c.db.Update(t); err != nil {
		return err
	}
	c.deleteCache(t.ID)
	return nil
}

func (c *CachedThreadRepo) DeleteByID(id uint) error {
	if err := c.db.DeleteByID(id); err != nil {
		return err
	}
	c.deleteCache(id)
	return nil
}

func (c *CachedThreadRepo) List(limit, offset int) ([]models.Thread, error) {
	return c.db.List(limit, offset)
}

func (c *CachedThreadRepo) ListAfter(cursorTime time.Time, cursorID uint, limit int) ([]models.Thread, error) {
	return c.db.ListAfter(cursorTime, cursorID, limit)
}

func (c *CachedThreadRepo) Count() (int64, error) {
	return c.db.Count()
}

func (c *CachedThreadRepo) ListByUserID(userID uint, limit, offset int) ([]models.Thread, error) {
	return c.db.ListByUserID(userID, limit, offset)
}

func (c *CachedThreadRepo) ListByUserIDAfter(userID uint, cursorTime time.Time, cursorID uint, limit int) ([]models.Thread, error) {
	return c.db.ListByUserIDAfter(userID, cursorTime, cursorID, limit)
}

func (c *CachedThreadRepo) CountByUserID(userID uint) (int64, error) {
	return c.db.CountByUserID(userID)
}

func (c *CachedThreadRepo) IncrementLikeCount(threadID uint, delta int) error {
	return c.db.IncrementLikeCount(threadID, delta)
}

func (c *CachedThreadRepo) GetLikeCount(threadID uint) (int64, error) {
	return c.db.GetLikeCount(threadID)
}

func (c *CachedThreadRepo) Transaction(fn func(tx *gorm.DB) error) error {
	txer, ok := c.db.(Transactioner)
	if !ok {
		return fmt.Errorf("不支持事务")
	}
	return txer.Transaction(fn)
}

func (c *CachedThreadRepo) WithTx(tx *gorm.DB) ThreadRepository {
	tr, ok := c.db.(ThreadRepoWithTx)
	if !ok {
		return c
	}
	return tr.WithTx(tx)
}
