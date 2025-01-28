package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"ws-trial/config"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func StartContainer(ctx context.Context, cli *client.Client) (string, error) {
	// Container configuration
	tmpDir, err := os.MkdirTemp("", "cpp-executor-*")
	if err := os.Chmod(tmpDir, 0777); err != nil {
		log.Fatalf("Failed to change permissions of temporary directory: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: config.IMAGE_NAME,
			Tty:   true, // Enable TTY for interactive input
			// User:  "1000:1000",
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:    512 * 1024 * 1024,
				CPUQuota:  100000,
				PidsLimit: &config.MAX_PROCESSES,
			},
			ReadonlyRootfs: true,
			NetworkMode:    "none",
			Binds: []string{
				tmpDir + ":/tmp/cpp", // Mount the temporary directory as writable
				tmpDir + ":/tmp:z",
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
