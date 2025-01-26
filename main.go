package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	imageName = "lightweight-code-execution-environment-2" // Replace with your Docker image name
)

func main() {
	// Fixed C++ code
	cppCode := `
#include <iostream>
#include <string>
#include <vector>
int main() {
    std::string input;
	std::vector<int> s(10, 0);
	for(auto &a: s) std::cin >> a;
    // std::getline(std::cin, input);
    for(auto &a: s) std::cout << a << " ";
    return 0;
}
`

	// Fixed stdin input
	stdinInput := `1
	2
	3
	4
	5
	6
	7
	8
	9
	10`

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	ctx := context.Background()

	containerID, err := startContainer(ctx, cli)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}
	log.Println("Container started with ID:", containerID)

	// Step 2: Run the C++ code in the container
	output, err := runCppInContainer(ctx, cli, containerID, cppCode, stdinInput)
	if err != nil {
		log.Fatalf("Failed to execute C++ code: %v", err)
	}
	fmt.Println("Output:")
	fmt.Println(output)

	// Step 3: Stop and remove the container
	log.Println("Stopping and removing container...")
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		log.Printf("Failed to stop container: %v", err)
	}
	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{}); err != nil {
		log.Printf("Failed to remove container: %v", err)
	}
}

var processCount int64 = 120

// startContainer starts a new container and returns its ID
func startContainer(ctx context.Context, cli *client.Client) (string, error) {
	// Container configuration
	tmpDir, err := os.MkdirTemp("", "cpp-executor-*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			Tty:   true, // Enable TTY for interactive input
			// User:  "1000:1000",
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:    512 * 1024 * 1024,
				CPUQuota:  100000,
				PidsLimit: &processCount,
			},
			ReadonlyRootfs: true,
			NetworkMode:    "none",
			Binds: []string{
				tmpDir + ":/tmp/cpp", // Mount the temporary directory as writable
			},
		}, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %v", err)
	}

	return resp.ID, nil
}

// runCppInContainer takes C++ code, creates a .cpp file, compiles, and runs it in a Docker container.
func runCppInContainer(ctx context.Context, cli *client.Client, containerID string, cppCode string, stdinInput string) (string, error) {
	// Step 1: Create a temporary .cpp file inside the container
	mkdirCmd := "mkdir -p /tmp/cpp"
	_, err := execInContainer(ctx, cli, containerID, mkdirCmd)
	if err != nil {
		return "", fmt.Errorf("failed to create /tmp/cpp directory: %v", err)
	}

	cppFileName := "/tmp/cpp/program.cpp"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(cppCode, "'", "'\\''"), cppFileName)
	_, err = execInContainer(ctx, cli, containerID, createFileCmd)
	if err != nil {
		return "", fmt.Errorf("failed to create .cpp file: %v", err)
	}

	// Step 2: Compile the .cpp file using g++
	compileCmd := fmt.Sprintf("g++ %s -o /tmp/cpp/program", cppFileName)
	compileOutput, err := execInContainer(ctx, cli, containerID, compileCmd)
	if err != nil {
		return "", fmt.Errorf("failed to compile C++ code: %v\nOutput: %s", err, compileOutput)
	}

	// Step 3: Run the compiled program with stdin input
	runCmd := fmt.Sprintf("echo '%s' | /tmp/cpp/program", stdinInput)
	runOutput, err := execInContainer(ctx, cli, containerID, runCmd)
	if err != nil {
		return "", fmt.Errorf("failed to run C++ program: %v\nOutput: %s", err, runOutput)
	}

	// Step 4: Clean up the compiled program and .cpp file
	cleanupCmd := fmt.Sprintf("rm %s /tmp/cpp/program", cppFileName)
	_, err = execInContainer(ctx, cli, containerID, cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}

	return runOutput, nil

}

func execInContainer(ctx context.Context, cli *client.Client, containerID string, command string) (string, error) {
	execConfig := types.ExecConfig{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
	}
	execID, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %v", err)
	}

	resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %v", err)
	}
	defer resp.Close()

	// Read the output
	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %v", err)
	}

	return string(output), nil
}
