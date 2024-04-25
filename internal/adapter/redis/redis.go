package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	SetCode(userID string, code string) error
	IsCodeValid(userID string, inputCode int) bool
}

type redisDB struct {
	client *redis.Client
}

func NewRedis(ipAddr string) (RedisInterface, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:7777",
		Password: "",
		DB:       0,
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

func (rdb redisDB) IsCodeValid(userID string, inputCode int) bool {
	res := rdb.client.Get(context.Background(), userID)
	code, _ := res.Int()
	return code == inputCode
}
