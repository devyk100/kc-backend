package worker

import (
	"fmt"
	"log"
	"strings"
)

func (w *Worker) createGoFile(code string) string {
	permCmd := "ls -ld /tmp/go"
	permOutput, err := w.dockerContainer.ExecInContainer(permCmd)
	if err != nil {

	}
	log.Println("Permissions of /tmp/go:", permOutput)

	goFileName := "/tmp/cpp/program.go"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(code, "'", "'\\''"), goFileName)
	_, err = w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {

	}
	return goFileName
}

func (w *Worker) compileGo(filename string) (string, error) {
	compileCmd := fmt.Sprintf("GOCACHE=/tmp/go/go-cache go build -o /tmp/go/program %s", filename)
	compileOutput, err := w.dockerContainer.ExecInContainer(compileCmd)
	if compileOutput != "" {
		err = fmt.Errorf("compile error")
	}
	return compileOutput, err
}

func (w *Worker) execGo(testcaseInput string) chan string {
	c := make(chan string)
	go func() {
		runCmd := fmt.Sprintf("echo '%s' | /tmp/go/program", testcaseInput)
		runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
		if err != nil {
			c <- err.Error()
		}
		c <- runOutput
	}()
	return c
}

func (w *Worker) cleanUpGo(filename string) {
	cleanupCmd := fmt.Sprintf("rm %s /tmp/go/program", filename)
	_, err := w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
