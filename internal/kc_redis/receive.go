package kcredis

func (r *RedisClient) Receive() (string, error) {
	result, err := r.redisClient.BLPop(r.ctx, 0, "jobs").Result()
	return result[1], err
}
