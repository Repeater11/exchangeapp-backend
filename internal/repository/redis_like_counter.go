package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	likeCountTTL       = 60 * time.Second
	likeCountTTLJitter = 10 * time.Second
	likeCountLockTTL   = 3 * time.Second
)

const incrIfExistsScript = `
if redis.call("EXISTS", KEYS[1]) == 1 then
	return redis.call("INCRBY", KEYS[1], ARGV[1])
end
return nil
`

const unlockIfMatchScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`

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
	_, err := c.rdb.Eval(context.Background(), incrIfExistsScript, []string{c.key(threadID)}, delta).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
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
	jitter := time.Duration(time.Now().UnixNano() % int64(likeCountTTLJitter))
	ttl := likeCountTTL + jitter
	return c.rdb.Set(context.Background(), c.key(threadID), value, ttl).Err()
}

func (c *RedisLikeCounter) lockKey(threadID uint) string {
	return fmt.Sprintf("thread:like:lock%d", threadID)
}

func (c *RedisLikeCounter) TryLockLikeCount(threadID uint, token string, ttl time.Duration) (bool, error) {
	return c.rdb.SetNX(context.Background(), c.lockKey(threadID), token, ttl).Result()
}

func (c *RedisLikeCounter) UnlockLikeCount(threadID uint, token string) error {
	_, err := c.rdb.Eval(context.Background(), unlockIfMatchScript, []string{c.lockKey(threadID)}, token).Result()
	return err
}
