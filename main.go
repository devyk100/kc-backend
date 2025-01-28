package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	rce "ws-trial/cmd/remote-code-execution"
	"ws-trial/config"
	"ws-trial/config/state"
	"ws-trial/config/types"
)

func main() {
	ctx := context.Background()
	// sqsClient, err := config.SqsClient()
	redisClient, err := config.RedisClient(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			payload := types.Payload_t{}
			val, err := redisClient.BLPop(ctx, 0, "jobs").Result()
			if err != nil {
				fmt.Println(err.Error(), "at unmarshall")
			}
			err = json.Unmarshal([]byte(val[1]), &payload)
			if err != nil {
				fmt.Println(err.Error(), "at unmarshall")
			}
			if !config.Running {
				return
			}

			// for _, val := range messages {
			// 	err := json.Unmarshal([]byte(*val.Body), &payload)
			// 	if err != nil {
			// 		fmt.Println(err.Error())
			// 	}

			// 	if !config.Running {
			// 		return
			// 	}

			// 	fmt.Println("key", payload.Key, "lang", payload.Language, "qid", payload.QuestionId)
			// 	payload.ReceiptHandle = *val.ReceiptHandle
			state.Jobs.Mut.Lock()
			state.Jobs.List = append(state.Jobs.List, payload)
			state.Jobs.Mut.Unlock()
			// }

		}
	}()
	go rce.CreateWorkers()
	<-sigChan
	config.Running = false
	state.Jobs.Mut.RLock()
	for i := 0; i < len(state.Jobs.List); i++ {
		str, err := json.Marshal(state.Jobs.List[i])
		if err != nil {
			fmt.Println(err.Error())
		}
		redisClient.RPush(ctx, "jobs", str)
	}
	state.Jobs.Mut.RUnlock()
	// 	// Fixed C++ code
	// 	cppCode := `
	// #include <iostream>
	// #include <string>
	// #include <vector>
	// int main() {
	//     std::string input;
	// 	std::vector<int> s(10, 0);
	// 	for(auto &a: s) std::cin >> a;
	//     // std::getline(std::cin, input);
	//     for(auto &a: s) std::cout << a << " ";
	//     return 0;
	// }
	// `

	// 	// Fixed stdin input
	// 	stdinInput := `1
	// 	2
	// 	3
	// 	4
	// 	5
	// 	6
	// 	7
	// 	8
	// 	9
	// 	10`

	// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// 	if err != nil {
	// 		log.Fatalf("Failed to create Docker client: %v", err)
	// 	}

	// 	ctx := context.Background()

	// 	containerID, err := startContainer(ctx, cli)
	// 	if err != nil {
	// 		log.Fatalf("Failed to start container: %v", err)
	// 	}
	// 	log.Println("Container started with ID:", containerID)

	// 	// Step 2: Run the C++ code in the container
	// 	output, err := runCppInContainer(ctx, cli, containerID, cppCode, stdinInput)
	// 	if err != nil {
	// 		log.Fatalf("Failed to execute C++ code: %v", err)
	// 	}
	// 	fmt.Println("Output:")
	// 	fmt.Println(output)

	// // Step 3: Stop and remove the container
	// log.Println("Stopping and removing container...")
	//
	//	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
	//		log.Printf("Failed to stop container: %v", err)
	//	}
	//
	//	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{}); err != nil {
	//		log.Printf("Failed to remove container: %v", err)
	//	}
}
