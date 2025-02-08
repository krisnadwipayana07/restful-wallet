package infrastructure

import (
	"context"

	"github.com/krisnadwipayana07/restful-fintech/configs"
	"github.com/redis/go-redis/v9"
)

// Create a Redis client
var ctx = context.Background()
var redisClient *redis.Client

func InitRedisConnection(config *configs.Config) (*redis.Client, error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
