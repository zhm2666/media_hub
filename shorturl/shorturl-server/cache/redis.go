package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)
import pkg_redis "shorturl/pkg/db/redis"

type redisKVCache struct {
	redisClient *redis.Client
	destroy     func()
}

func newRedisKVCache(client *redis.Client, destroy func()) KVCache {
	return &redisKVCache{
		redisClient: client,
		destroy:     destroy,
	}
}

func getKey(key string) string {
	return pkg_redis.GetKey(key)
}
func (c *redisKVCache) Get(key string) (string, error) {
	key = getKey(key)
	res, err := c.redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}
func (c *redisKVCache) Set(key, value string, ttl int) error {
	key = getKey(key)
	return c.redisClient.SetEx(context.Background(), key, value, time.Second*time.Duration(ttl)).Err()
}
func (c *redisKVCache) Destroy() {
	if c.destroy != nil {
		c.destroy()
	}
}

// redis 工厂
type redisCacheFactory struct {
	redisPool pkg_redis.RedisPool
}

func NewRedisCacheFactory(redisPool pkg_redis.RedisPool) CacheFactory {
	return &redisCacheFactory{
		redisPool: redisPool,
	}
}

func (f *redisCacheFactory) NewKVCache() KVCache {
	client := f.redisPool.Get()
	return newRedisKVCache(client, func() {
		f.redisPool.Put(client)
	})
}
