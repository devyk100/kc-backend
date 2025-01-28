package rce

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"
	"ws-trial/cmd/docker"
	"ws-trial/config"
	"ws-trial/config/state"
	"ws-trial/config/types"

	"github.com/docker/docker/client"
)

var RunningWorkerCount int32 = 0

func CreateWorkers() {
	for {
		if !config.Running {
			return
		}
		for atomic.LoadInt32(&RunningWorkerCount) < int32(config.MAX_CONTAINERS) {
			if !config.Running {
				return
			}
			go RunWorker(int(atomic.LoadInt32(&RunningWorkerCount)))
			atomic.AddInt32(&RunningWorkerCount, 1)
		}
	}
}

func RunWorker(workerId int) {
	ctx := context.TODO()
	fmt.Println("A new container spawned")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	containerId, err := docker.StartContainer(ctx, cli)
	defer func() {
		fmt.Println("THis is the defer of this function")
		atomic.AddInt32(&RunningWorkerCount, -1)
		err := cli.ContainerKill(context.Background(), containerId, "SIGKILL")
		if err != nil {
			log.Fatal("Error inspecting container: ", err)
		}

	}()
	if err != nil {
		fmt.Println("Error", err.Error())
		return
	}
	val := false
	go func() {
		time.Sleep(time.Second * 10)
		val = true
	}()
	for {
		if val {
			return
		}
		if !config.Running {
			return
		}
		if len(state.Jobs.List) == 0 {
			// fmt.Println("No work to do", workerId)
			continue
		}
		var Job *types.Payload_t
		state.Jobs.Mut.Lock()
		if len(state.Jobs.List) != 0 {
			JobVal := state.Jobs.List[0]
			Job = &JobVal
			fmt.Println("Got this job", Job.Key, "for the worker", workerId)
			fmt.Println("deleted that job")
			state.Jobs.List = state.Jobs.List[1:]
		}
		state.Jobs.Mut.Unlock()

		if Job == nil {
			continue
		}
	}

}
