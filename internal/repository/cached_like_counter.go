package repository

import "gorm.io/gorm"

type threadLikeCache interface {
	ThreadLikeCounter
	setLikeCount(threadID uint, value int64) error
}

type CachedThreadLikeCounter struct {
	db    ThreadRepository
	cache threadLikeCache
}

func NewCachedThreadLikeCounter(db ThreadRepository, cache *RedisLikeCounter) *CachedThreadLikeCounter {
	return &CachedThreadLikeCounter{
		db:    db,
		cache: cache,
	}
}

func (c *CachedThreadLikeCounter) IncrementLikeCount(threadID uint, delta int) error {
	if err := c.db.IncrementLikeCount(threadID, delta); err != nil {
		return err
	}

	_ = c.cache.IncrementLikeCount(threadID, delta)
	return nil
}

func (c *CachedThreadLikeCounter) GetLikeCount(threadID uint) (int64, error) {
	val, err := c.cache.GetLikeCount(threadID)
	if err == nil {
		return val, nil
	}

	val, err = c.db.GetLikeCount(threadID)
	if err != nil {
		return 0, err
	}

	_ = c.cache.setLikeCount(threadID, val)
	return val, nil
}

func (c *CachedThreadLikeCounter) WithTx(tx *gorm.DB) ThreadLikeCounter {
	tr, ok := c.db.(ThreadRepoWithTx)
	if !ok {
		return c
	}
	return &CachedThreadLikeCounter{
		db:    tr.WithTx(tx),
		cache: c.cache,
	}
}
