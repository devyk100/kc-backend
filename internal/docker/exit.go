package docker

import (
	"fmt"
	"log"
)

func (d *Docker) Exit() {
	fmt.Println("Exiting the docker container")
	err := d.cli.ContainerKill(d.ctx, d.containerId, "SIGKILL")
	if err != nil {
		log.Fatal("Error inspecting container: ", err)
	}
}
