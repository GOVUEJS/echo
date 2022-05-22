package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"myapp/util"
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
		Addr:     *util.Host + ":6379",
		Password: *util.Password, // no password set
		DB:       0,              // use default DB
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
