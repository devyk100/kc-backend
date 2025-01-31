package main

import (
	"context"
	"ws-trial/internal/orchestrator"
)

func main() {
	ctx := context.Background()
	o := orchestrator.Orchestrator{}
	o.Run(ctx)
}
