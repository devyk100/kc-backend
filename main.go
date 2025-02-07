package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	ws_server "ws-trial/cmd/yjs-ws-server"
	"ws-trial/internal/judge_orchestrator"
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

	go func() {
		ws_server.Start()
	}()

	o := judge_orchestrator.Orchestrator{}
	o.Run(ctx)
}
