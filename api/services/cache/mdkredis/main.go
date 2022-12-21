package mdkredis

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mesirendon/gredis/services/cache"
)

type Cache struct {
	Instance redis.Cmdable
	Context  context.Context
}

// Creates a new instance of the cache service for the given
// Redis connection.
//
// Example:
//
//	client := redis.NewClient(&redis.Options{
//		Addr: mr.Addr(),
//	})
func New(i redis.Cmdable) cache.ICache {
	return &Cache{
		Instance: i,
		Context:  context.Background(),
	}
}

func (c *Cache) Set(key string, value any, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.Instance.Set(c.Context, key, jsonValue, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Get(key string, output any) error {
	val, err := c.Instance.Get(c.Context, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), &output); err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetRecordsByPattern(regexp string, output any) error {
	var cursor uint64
	var allKeys = make([]string, 0)
	var values = make([]string, 0)

	iterator := c.Instance.Scan(c.Context, cursor, regexp+"*", 0).Iterator()

	for iterator.Next(c.Context) {
		allKeys = append(allKeys, iterator.Val())
	}

	if err := iterator.Err(); err != nil {
		return err
	}

	for _, key := range allKeys {
		val, _ := c.Instance.Get(c.Context, key).Result()
		values = append(values, val)
	}

	if err := json.Unmarshal([]byte("["+strings.Join(values, ",")+"]"), &output); err != nil {
		return err
	}

	return nil
}

func (c *Cache) TTL(key string) (time.Duration, error) {
	duration := c.Instance.TTL(c.Context, key)

	return duration.Val(), nil
}

func (c *Cache) KeyPatternRecordSize(regexp string) (int, error) {
	var cursor uint64
	var size int

	for {
		keys, cursor, err := c.Instance.Scan(c.Context, cursor, regexp, 10).Result()
		if err != nil {
			return 0, err
		}

		size += len(keys)

		if cursor == 0 {
			break
		}
	}

	return size, nil
}

func (c *Cache) Close() error {
	client := c.Instance.(*redis.Client)
	if err := client.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) FlushDB() error {
	if err := c.Instance.FlushDB(c.Context).Err(); err != nil {
		return err
	}

	return nil
}
