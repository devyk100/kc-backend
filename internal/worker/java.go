package worker

import (
	"fmt"
	"log"
	"strings"
)

func (w *Worker) runJavaInContainer(job *Job) {
	permCmd := "ls -ld /tmp/java"
	permOutput, err := w.dockerContainer.ExecInContainer(permCmd)
	if err != nil {

	}
	log.Println("Permissions of /tmp/java:", permOutput)

	javaFileName := "/tmp/cpp/Program.java"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(job.Code, "'", "'\\''"), javaFileName)
	_, err = w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {

	}

	compileCmd := fmt.Sprintf("javac %s", javaFileName)
	compileOutput, err := w.dockerContainer.ExecInContainer(compileCmd)
	fmt.Println(compileOutput)

	if err != nil {
	}
	// fetch all the input test cases
	// fetch all the required outputs

	stdinInput := "1"
	runCmd := fmt.Sprintf("echo '%s' | java -cp /tmp/java Program", stdinInput)
	runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
	fmt.Print(runOutput)
	if err != nil {
		// Handle error
	}

	cleanupCmd := fmt.Sprintf("rm %s /tmp/java/Program.class", javaFileName)
	_, err = w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}

}
