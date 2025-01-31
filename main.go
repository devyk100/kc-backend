package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"ws-trial/internal/orchestrator"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Print("Called")
		cancel()
	}()

	o := orchestrator.Orchestrator{}
	o.Run(ctx)
}
