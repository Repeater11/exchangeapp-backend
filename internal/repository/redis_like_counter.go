package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisLikeCounter struct {
	rdb *redis.Client
}

func NewRedisLikeCounter(rdb *redis.Client) *RedisLikeCounter {
	return &RedisLikeCounter{
		rdb: rdb,
	}
}

func (c *RedisLikeCounter) key(threadID uint) string {
	return fmt.Sprintf("thread:like:%d", threadID)
}

func (c *RedisLikeCounter) IncrementLikeCount(threadID uint, delta int) error {
	return c.rdb.IncrBy(context.Background(), c.key(threadID), int64(delta)).Err()
}

func (c *RedisLikeCounter) GetLikeCount(threadID uint) (int64, error) {
	val, err := c.rdb.Get(context.Background(), c.key(threadID)).Int64()
	if err == redis.Nil {
		return 0, ErrLikeCountNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("获取点赞数失败：%w", err)
	}
	return val, nil
}

func (c *RedisLikeCounter) setLikeCount(threadID uint, value int64) error {
	return c.rdb.Set(context.Background(), c.key(threadID), value, 0).Err()
}
