package orchestrator

import (
	"context"
	kcredis "ws-trial/internal/kc_redis"
	"ws-trial/internal/worker"
)

type Orchestrator struct {
	redisClient *kcredis.RedisClient
	ctx         context.Context
	Workers     []worker.Worker
}

func (o *Orchestrator) Run(ctx context.Context) {
	o.ctx = ctx
	redisClient, err := kcredis.CreateRedisClient(ctx)
	if err != nil {
		return
	}
	o.redisClient = redisClient
	o.Workers = make([]worker.Worker, 3)

	for i := 0; i < len(o.Workers); i++ {
		ctx, cancel := context.WithCancel(o.ctx)
		o.Workers[i].Run(ctx, cancel)
	}

	defer func() {
		for i := 0; i < len(o.Workers); i++ {
			o.Workers[i].Cancel()
		}
	}()

	for {
		select {
		case <-o.ctx.Done():
			return
		default:
			continue
		}
	}

}
