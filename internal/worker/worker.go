package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"ws-trial/internal/docker"
	kcredis "ws-trial/internal/kc_redis"
)

type Job struct {
	Code     string `json:"code"`
	Qid      int    `json:"qid"`
	Language string `json:"lang"`
	QueryKey string `json:"querykey"`
}

type FinishedPayload struct {
	Message   string        `json:"message"`
	Where     string        `json:"where"`
	TimeTaken time.Duration `json:"timetaken"`
}

type Worker struct {
	redisClient     *kcredis.RedisClient
	ctx             context.Context
	Cancel          context.CancelFunc
	dockerContainer docker.Docker
}

func (w *Worker) Exit() {
	w.dockerContainer.Exit()
	w.redisClient.Exit()
}

func (w *Worker) Run(ctx context.Context, cancel context.CancelFunc) {
	var payload Job
	var val string

	w.ctx = ctx
	w.Cancel = cancel
	w.dockerContainer.StartContainer(context.Background())
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
				w.redisClient.RePut(val)
				return
			default:
				{
					val, err := w.redisClient.Receive()
					if val == "" {
						continue
					}

					if err != nil {
						w.Exit()
					}
					err = json.Unmarshal([]byte(val), &payload)
					if err != nil {
						w.Exit()
						return
					}
					fmt.Println("Got this", val)
					resp := w.Exec(payload)
					fmt.Print(resp.Message, resp.TimeTaken, resp.Where)
					respStr, err := json.Marshal(resp)
					if err != nil {
						fmt.Println("Error at marshalling", err.Error())
					}
					w.redisClient.PutFinishedJob(payload.QueryKey, string(respStr))
				}
			}
		}
	}()
}
