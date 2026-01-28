package repository

import (
	"strconv"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type threadLikeCache interface {
	ThreadLikeCounter
	setLikeCount(threadID uint, value int64) error
	TryLockLikeCount(threadID uint, token string, ttl time.Duration) (bool, error)
	UnlockLikeCount(threadID uint, token string) error
	MarkDirty(threadID uint) error
}

type CachedThreadLikeCounter struct {
	db    ThreadRepository
	cache threadLikeCache
	sf    *singleflight.Group
}

func NewCachedThreadLikeCounter(db ThreadRepository, cache *RedisLikeCounter) *CachedThreadLikeCounter {
	return &CachedThreadLikeCounter{
		db:    db,
		cache: cache,
		sf:    &singleflight.Group{},
	}
}

func (c *CachedThreadLikeCounter) IncrementLikeCount(threadID uint, delta int) error {
	if err := c.cache.IncrementLikeCount(threadID, delta); err != nil {
		return err
	}

	return c.cache.MarkDirty(threadID)
}

func (c *CachedThreadLikeCounter) GetLikeCount(threadID uint) (int64, error) {
	if val, err := c.cache.GetLikeCount(threadID); err == nil {
		return val, nil
	}

	key := "thread_like_count:" + strconv.FormatUint(uint64(threadID), 10)
	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		token := strconv.FormatInt(time.Now().UnixNano(), 10)
		locked, _ := c.cache.TryLockLikeCount(threadID, token, likeCountLockTTL)
		if locked {
			defer c.cache.UnlockLikeCount(threadID, token)
			val, err := c.db.GetLikeCount(threadID)
			if err != nil {
				return nil, err
			}
			_ = c.cache.setLikeCount(threadID, val)
			return val, err
		}

		time.Sleep(20 * time.Millisecond)
		if val, err := c.cache.GetLikeCount(threadID); err == nil {
			return val, nil
		}

		val, err := c.db.GetLikeCount(threadID)
		if err != nil {
			return nil, err
		}
		_ = c.cache.setLikeCount(threadID, val)
		return val, nil
	})
	if err != nil {
		return 0, err
	}
	return v.(int64), nil
}

func (c *CachedThreadLikeCounter) WithTx(tx *gorm.DB) ThreadLikeCounter {
	tr, ok := c.db.(ThreadRepoWithTx)
	if !ok {
		return c
	}
	return &CachedThreadLikeCounter{
		db:    tr.WithTx(tx),
		cache: c.cache,
		sf:    c.sf,
	}
}
