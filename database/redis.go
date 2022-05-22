package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"myapp/config"
	"time"
)

var (
	redisContext = context.Background()
	redisClient  *redis.Client
)

func init() {
	InitRedis()
}

func InitRedis() {
	initRedisClient()
	initRedisContext()
}

func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     *config.Host + ":6379",
		Password: *config.Password, // no password set
		DB:       0,                // use default DB
	})
}

func initRedisContext() {
	redisContext = context.Background()
}

func GetRedis() *redis.Client {
	if redisClient == nil {
		InitRedis()
	}
	return redisClient
}

func GetRedisContext() *context.Context {
	if redisContext == nil {
		initRedisContext()
	}
	return &redisContext
}

func SetRedisSession(sessionId string, refreshToken *string, accessToken *string) error {
	ctx := *GetRedisContext()
	if _, err := redisClient.Pipelined(ctx, func(redis redis.Pipeliner) error {
		redis.HSet(ctx, sessionId, "refreshToken", refreshToken)
		redis.HSet(ctx, sessionId, "accessToken", accessToken)
		redis.Expire(ctx, sessionId, time.Hour)
		return nil
	}); err != nil {
		return err
	}
	return nil
}
