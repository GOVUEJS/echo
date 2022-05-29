package session

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"myapp/config"
	"myapp/model"
)

var (
	redisContext = context.Background()
	redisClient  *redis.Client
)

func InitRedis() error {
	err := initRedisClient()
	if err != nil {
		panic(err)
	}
	initRedisContext()
	return nil
}

func initRedisClient() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Config.Redis.Ip, config.Config.Redis.Port),
		Password: config.Config.Redis.Password, // no password set
		DB:       0,                            // use default DB
	})

	return redisClient.Ping(redisContext).Err()
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
