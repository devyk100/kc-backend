package docker

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-units"
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
			Env: []string{
				"LANG=en_US.UTF-8",
				"LC_ALL=en_US.UTF-8",
				"DONT_POLLUTE_OUTPUT_WITH_UTF8=1", // If you need it
			},

			// User:  "1000:1000",
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:    4 * 512 * 1024 * 1024,
				CPUQuota:  200000,
				PidsLimit: &MAX_PROCESSES,
				Ulimits: []*units.Ulimit{
					{
						Name: "nofile",
						Hard: 65535,
						Soft: 65535,
					},
				},
			},
			AutoRemove:     true,
			ReadonlyRootfs: true,
			NetworkMode:    "none",
			Binds: []string{
				tmpDir + ":/tmp/cpp", // Mount the temporary directory as writable
				tmpDir + ":/tmp:z",
				tmpDir + ":/tmp/go",
				tmpDir + ":/tmp/java",
				tmpDir + ":/tmp/python",
				tmpDir + ":/tmp/js",
			},
		}, nil, nil, "")
	if err != nil {
		return err
	}

	if err := d.cli.ContainerStart(d.ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	d.containerId = resp.ID
	return nil
}
