package main

import (
	"context"
	"fmt"
	"ws-trial/internal/docker"
)

func main() {
	ctx := context.Background()
	d := docker.Docker{}
	err := d.StartContainer(ctx)
	permCmd := "ls -ld /tmp/cpp"
	permOutput, err := d.ExecInContainer(permCmd)
	fmt.Println(permOutput)
	d.Exit()
	if err != nil {
		return
	}
}
