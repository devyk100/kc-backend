package judge_orchestrator

import (
	"context"
	"fmt"
	"ws-trial/internal/worker"
)

type Orchestrator struct {
	ctx     context.Context
	Workers []worker.Worker
}

func (o *Orchestrator) Run(ctx context.Context) {
	o.ctx = ctx
	o.Workers = make([]worker.Worker, 1)
	fmt.Println("Starting ", len(o.Workers), "Workers")
	for i := 0; i < len(o.Workers); i++ {
		ctx, cancel := context.WithCancel(o.ctx)
		o.Workers[i].Run(ctx, cancel)
	}

	for {
		select {
		case <-o.ctx.Done():
			for i := 0; i < len(o.Workers); i++ {
				o.Workers[i].Exit()
				o.Workers[i].Cancel()
			}
			return
		default:
			continue
		}
	}

}
