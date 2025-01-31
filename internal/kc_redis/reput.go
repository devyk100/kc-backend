package kcredis

import "fmt"

func (r *RedisClient) RePut(val string) {
	_, err := r.redisClient.LPush(r.ctx, "jobs").Result()
	if err != nil {
		fmt.Println(err.Error())
	}
}
