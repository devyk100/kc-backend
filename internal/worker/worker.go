package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"ws-trial/internal/docker"
	kcredis "ws-trial/internal/kc_redis"
)

type Job struct {
	Code     string `json:"code"`
	Qid      int    `json:"qid"`
	Language string `json:"lang"`
	QueryKey string `json:"querykey"`
}

type Worker struct {
	redisClient     *kcredis.RedisClient
	ctx             context.Context
	Cancel          context.CancelFunc
	dockerContainer docker.Docker
}

func (w *Worker) Exit() {
	w.redisClient.Exit()
}

func (w *Worker) Run(ctx context.Context, cancel context.CancelFunc) {
	w.ctx = ctx
	w.Cancel = cancel
	w.dockerContainer.StartContainer(w.ctx)
	redisClient, err := kcredis.CreateRedisClient(w.ctx)
	if err != nil {
		w.Exit()
		return
	}
	w.redisClient = redisClient
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				fmt.Println("Worker received stop signal. Exiting...")
				return
			default:
				{
					var payload Job
					val, err := w.redisClient.Receive()
					if err != nil {
						w.Exit()
					}
					err = json.Unmarshal([]byte(val), &payload)
					if err != nil {
						w.Exit()
						return
					}
					fmt.Println("Got this", val)
					w.Exec(payload)
				}
			}
		}
	}()
}
