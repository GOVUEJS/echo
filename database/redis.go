package database

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"myapp/util"
	"time"
)

var (
	ctx   = context.Background()
	cache = new(Redis)
)

func init() {
	InitRedis()
}

func InitRedis() {
	cache.redisClient = redis.NewClient(&redis.Options{
		Addr:     *util.Host + ":6379",
		Password: *util.Password, // no password set
		DB:       0,              // use default DB
	})
}

func GetRedis() *Redis {
	if cache == nil {
		InitRedis()
	}
	return cache
}

type Redis struct {
	redisClient *redis.Client
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) {
	err := r.redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
}

func (r *Redis) Get(key string) (value string, err error) {
	value, err = r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("key does not exist")
	} else if err != nil {
		return "", err
	}
	return
}
