package worker

import (
	"fmt"
	"log"
	"strings"
)

func (w *Worker) runJavaScriptInContainer(job *Job) {
	jsFileName := "/tmp/js/program.js"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(job.Code, "'", "'\\''"), jsFileName)
	_, err := w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {
		// Handle error
	}

	stdinInput := "1"
	runCmd := fmt.Sprintf("echo '%s' | node %s", stdinInput, jsFileName)
	runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
	fmt.Print(runOutput)
	if err != nil {
		// Handle error
	}

	cleanupCmd := fmt.Sprintf("rm %s", jsFileName)
	_, err = w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
