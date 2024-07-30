package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"os"
	"time"
)

type CacheStruct struct {
	ctx    context.Context
	client *redis.Client
}

var Cache *CacheStruct

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPass := os.Getenv("REDIS_PASS")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPass,
		DB:       6,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		Logger.Error("Failed to connect to Redis: ", zap.Error(err))
		panic(err)
		return
	}

	Cache = &CacheStruct{
		ctx:    context.Background(),
		client: redisClient,
	}
}

func (cache *CacheStruct) Set(key string, value interface{}, ttl time.Duration) error {
	return cache.client.Set(cache.ctx, key, value, ttl).Err()
}

func (cache *CacheStruct) Get(key string) (string, error) {
	return cache.client.Get(cache.ctx, key).Result()
}
