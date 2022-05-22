package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"myapp/config"
	"myapp/model"
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

func SetRedisSession(sessionId string, redisSession *model.RedisSession) error {
	ctx := *GetRedisContext()
	if _, err := redisClient.Pipelined(ctx, func(redis redis.Pipeliner) error {
		redis.HSet(ctx, sessionId, "email", *redisSession.Email)
		redis.HSet(ctx, sessionId, "ip", *redisSession.Ip)
		redis.HSet(ctx, sessionId, "accessToken", *redisSession.AccessToken)
		redis.HSet(ctx, sessionId, "refreshToken", *redisSession.RefreshToken)
		redis.Expire(ctx, sessionId, time.Hour)
		return nil
	}); err != nil {
		return err
	}
	return nil
}
