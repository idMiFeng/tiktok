package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedisClient() error {
	ctx := context.Background()

	// 创建Redis客户端连接对象
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "124.71.74.243:6379",
		Password: "",
		DB:       0,
	})

	// 测试连接
	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return err
	}

	fmt.Println("Connected to Redis:", pong)
	return nil
}
