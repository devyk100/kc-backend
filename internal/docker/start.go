package docker

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func (d *Docker) StartContainer(ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	d.cli = cli
	d.ctx = ctx
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "executor-*")
	if err := os.Chmod(tmpDir, 0777); err != nil {
		log.Fatalf("Failed to change permissions of temporary directory: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}

	resp, err := d.cli.ContainerCreate(d.ctx,
		&container.Config{
			Image: IMAGE_NAME,
			Tty:   true,
			// User:  "1000:1000",
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:    512 * 1024 * 1024,
				CPUQuota:  100000,
				PidsLimit: &MAX_PROCESSES,
			},
			ReadonlyRootfs: true,
			NetworkMode:    "none",
			Binds: []string{
				tmpDir + ":/tmp/cpp", // Mount the temporary directory as writable
				tmpDir + ":/tmp:z",
			},
		}, nil, nil, "")
	if err != nil {
		return err
	}

	if err := d.cli.ContainerStart(d.ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	d.containerId = resp.ID
	//fmt.Println(d.cli, d.ctx, d.containerId)
	if err != nil {
		return err
	}
	return nil
}
