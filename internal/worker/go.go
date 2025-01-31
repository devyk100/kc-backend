package worker

import (
	"fmt"
	"log"
	"strings"
)

func (w *Worker) runGoInContainer(job *Job) {
	permCmd := "ls -ld /tmp/go"
	permOutput, err := w.dockerContainer.ExecInContainer(permCmd)
	if err != nil {

	}
	log.Println("Permissions of /tmp/go:", permOutput)

	goFileName := "/tmp/cpp/program.go"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(job.Code, "'", "'\\''"), goFileName)
	_, err = w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {

	}
	_, err = w.dockerContainer.ExecInContainer("mkdir -p /tmp/go/go-cache")
	if err != nil {
		w.Exit()
		fmt.Println("Error occurred and exiting the container at mkdir golang")
		return
	}

	compileCmd := fmt.Sprintf("GOCACHE=/tmp/go/go-cache go build -o /tmp/go/program %s", goFileName)
	compileOutput, err := w.dockerContainer.ExecInContainer(compileCmd)
	fmt.Println(compileOutput)
	if err != nil {
	}
	// fetch all the input test cases
	// fetch all the required outputs
	stdinInput := "1"
	runCmd := fmt.Sprintf("echo '%s' | /tmp/go/program", stdinInput)
	runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
	fmt.Print(runOutput)
	if err != nil {
		// return "", fmt.Errorf("failed to run C++ program: %v\nOutput: %s", err, runOutput)
	}

	cleanupCmd := fmt.Sprintf("rm %s /tmp/go/program", goFileName)
	_, err = w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
