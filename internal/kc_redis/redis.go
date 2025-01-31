package kcredis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	ctx         context.Context
	redisClient *redis.Client
}

func CreateRedisClient(ctx context.Context) (*RedisClient, error) {
	// fmt.Println(os.Getenv(""), env.REDIS_URL)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error", err.Error())
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:      os.Getenv("REDIS_URL"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		TLSConfig: &tls.Config{},
		DB:        0,
	})
	// ctx := context.Background()
	// Ping the Redis server to check the connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)
	RedisClient_v := RedisClient{redisClient: rdb, ctx: ctx}
	return &RedisClient_v, err
}

func (r *RedisClient) Exit() {
	r.redisClient.Close()
}
