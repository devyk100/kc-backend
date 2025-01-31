package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"ws-trial/db"
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
	Message   string `json:"message"`
	Where     string `json:"where"`
	TimeTaken int32  `json:"timetaken"`
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
	query, pool, err := db.InitDb(w.ctx)
	if err != nil {
		fmt.Println("Error occured at db initialisation ", err.Error())
	}
	w.dockerContainer.StartContainer(context.Background()) // Using this child context is making things complicated, and the docker container client detaches with the bottommost child context called cancel, and hence killing the containers is out of the question, rather let it rely on the parent context.
	redisClient, err := kcredis.CreateRedisClient(w.ctx)
	if err != nil {
		w.Exit()
		return
	}
	w.redisClient = redisClient
	go func() {
		defer pool.Close()
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
					resp := w.Exec(payload, query)
					fmt.Print(resp.Message, resp.TimeTaken, resp.Where)
					respStr, err := json.Marshal(resp)
					if err != nil {
						fmt.Println("Error at marshalling", err.Error())
					}
					w.redisClient.PutFinishedJob(payload.QueryKey, string(respStr))
					valN, err := query.InsertSubmission(w.ctx, db.InsertSubmissionParams{
						Code:       payload.Code,
						Message:    resp.Message,
						Correct:    resp.Message == "Correct",
						QuestionID: int32(payload.Qid),
						Language:   payload.Language,
						Duration:   int64(resp.TimeTaken),
					})
					if err != nil {
						fmt.Println("Err", err.Error())
					}
					fmt.Print(valN)
				}
			}
		}
	}()
}
