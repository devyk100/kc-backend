package docker

import (
	"bytes"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

func (d *Docker) ExecInContainer(command string) (string, error) {
	execConfig := container.ExecOptions{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
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

func (d *Docker) ExecInContainerStdCopy(command string) (string, error) {
	// Create an exec configuration
	execConfig := types.ExecConfig{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
	}

	// Create an exec instance in the container
	execID, err := d.cli.ContainerExecCreate(d.ctx, d.containerId, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %v", err)
	}

	// Attach to the exec instance
	resp, err := d.cli.ContainerExecAttach(d.ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %v", err)
	}
	defer resp.Close()

	// Use buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer

	// Use stdcopy.StdCopy to separate stdout and stderr
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to copy output: %v", err)
	}

	// Combine stdout and stderr (if stderr is not empty)
	output := stdoutBuf.String()
	if stderrBuf.Len() > 0 {
		output += "\n" + stderrBuf.String()
	}

	return output, nil
}
