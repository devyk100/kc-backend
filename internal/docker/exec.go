package docker

import (
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
)

func (d *Docker) ExecInContainer(command string) (string, error) {
	execConfig := container.ExecOptions{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
	}
	execID, err := d.cli.ContainerExecCreate(d.ctx, d.containerId, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %v", err)
	}

	resp, err := d.cli.ContainerExecAttach(d.ctx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %v", err)
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %v", err)
	}

	return string(output), nil
}
