package worker

import (
	"fmt"
	"log"
	"strings"
)

func (w *Worker) runPythonInContainer(job *Job) {
	pyFileName := "/tmp/python/program.py"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(job.Code, "'", "'\\''"), pyFileName)
	_, err := w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {
		// Handle error
	}

	stdinInput := "1"
	runCmd := fmt.Sprintf("echo '%s' | python3 %s", stdinInput, pyFileName)
	runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
	fmt.Print(runOutput)
	if err != nil {
		// Handle error
	}

	cleanupCmd := fmt.Sprintf("rm %s", pyFileName)
	_, err = w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
