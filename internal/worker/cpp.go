package worker

import (
	"fmt"
	"log"
	"strings"
)

func (w *Worker) runCppInContainer(job *Job) {
	permCmd := "ls -ld /tmp/cpp"
	permOutput, err := w.dockerContainer.ExecInContainer(permCmd)
	if err != nil {

	}
	log.Println("Permissions of /tmp/cpp:", permOutput)

	cppFileName := "/tmp/cpp/program.cpp"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(job.Code, "'", "'\\''"), cppFileName)
	_, err = w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {

	}

	compileCmd := fmt.Sprintf("g++ %s -o /tmp/cpp/program", cppFileName)
	compileOutput, err := w.dockerContainer.ExecInContainer(compileCmd)
	fmt.Println(compileOutput)
	if err != nil {
	}
	// fetch all the input test cases
	// fetch all the required outputs

	// runCmd := fmt.Sprintf("echo '%s' | /tmp/cpp/program", stdinInput)
	// runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to run C++ program: %v\nOutput: %s", err, runOutput)
	// }

	cleanupCmd := fmt.Sprintf("rm %s /tmp/cpp/program", cppFileName)
	_, err = w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
