package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	SetCode(userID string, code string) error
}

type redisDB struct {
	client *redis.Client
}

func NewRedis(ipAddr string) (RedisInterface, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: ipAddr,
		DB:   0,
	})
	status := rdb.Ping(context.Background())

	return redisDB{
		client: rdb,
	}, status.Err()
}

func (rdb redisDB) SetCode(userID string, code string) error {
	status := rdb.client.Set(context.Background(), code, userID, time.Second*30)
	return status.Err()
}
