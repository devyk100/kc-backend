package kcredis

import "time"

func (r *RedisClient) PutFinishedJob(key string, payload string) {
	r.redisClient.Set(r.ctx, key, payload, time.Minute*20)
}
