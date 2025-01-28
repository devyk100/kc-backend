package docker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/client"
)

func runCppInContainer(ctx context.Context, cli *client.Client, containerID string, cppCode string, stdinInput string) (string, error) {
	// Step 1: Create a temporary .cpp file inside the container
	permCmd := "ls -ld /tmp/cpp"
	permOutput, err := ExecInContainer(ctx, cli, containerID, permCmd)
	if err != nil {
		return "", fmt.Errorf("failed to check /tmp/cpp permissions: %v\nOutput: %s", err, permOutput)
	}
	log.Println("Permissions of /tmp/cpp:", permOutput)
	// mkdirCmd := "mkdir -p /tmp/cpp"
	// _, err = ExecInContainer(ctx, cli, containerID, mkdirCmd)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to create /tmp/cpp directory: %v", err)
	// }

	cppFileName := "/tmp/cpp/program.cpp"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(cppCode, "'", "'\\''"), cppFileName)
	_, err = ExecInContainer(ctx, cli, containerID, createFileCmd)
	if err != nil {
		return "", fmt.Errorf("failed to create .cpp file: %v", err)
	}
	// permCmd2 := "ls /tmp/cpp"
	// permOutput2, err := ExecInContainer(ctx, cli, containerID, permCmd2)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to check /tmp/cpp permissions: %v\nOutput: %s", err, permOutput2)
	// }
	// log.Println("Permissions of /tmp/cpp:", permOutput2)
	// Step 2: Compile the .cpp file using g++
	compileCmd := fmt.Sprintf("g++ %s -o /tmp/cpp/program", cppFileName)
	compileOutput, err := ExecInContainer(ctx, cli, containerID, compileCmd)
	if err != nil {
		return "", fmt.Errorf("failed to compile C++ code: %v\nOutput: %s", err, compileOutput)
	}

	// Step 3: Run the compiled program with stdin input
	runCmd := fmt.Sprintf("echo '%s' | /tmp/cpp/program", stdinInput)
	runOutput, err := ExecInContainer(ctx, cli, containerID, runCmd)
	if err != nil {
		return "", fmt.Errorf("failed to run C++ program: %v\nOutput: %s", err, runOutput)
	}

	// Step 4: Clean up the compiled program and .cpp file
	cleanupCmd := fmt.Sprintf("rm %s /tmp/cpp/program", cppFileName)
	_, err = ExecInContainer(ctx, cli, containerID, cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}

	return runOutput, nil

}
