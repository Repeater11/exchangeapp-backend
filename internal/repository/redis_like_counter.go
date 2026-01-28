package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	likeCountLockTTL = 3 * time.Second
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

const likeDirtyKey = "thread:like:dirty"

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

func (c *RedisLikeCounter) MarkDirty(threadID uint) error {
	return c.rdb.SAdd(context.Background(), likeDirtyKey, threadID).Err()
}

func (c *RedisLikeCounter) PopDirty(limit int) ([]uint, error) {
	vals, err := c.rdb.SPopN(context.Background(), likeDirtyKey, int64(limit)).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(vals))
	for _, v := range vals {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids, nil
}
