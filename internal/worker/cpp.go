package worker

import (
	"fmt"
	"log"
)

func (w *Worker) createCppFile(code string) string {
	permCmd := "ls -ld /tmp/cpp"
	permOutput, err := w.dockerContainer.ExecInContainer(permCmd)
	if err != nil {

	}
	log.Println("Permissions of /tmp/cpp:", permOutput)

	cppFileName := "/tmp/cpp/program.cpp"
	createFileCmd := fmt.Sprintf("echo '%s' > %s", code, cppFileName)
	_, err = w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {

	}
	return cppFileName
}

func (w *Worker) compileCpp(filename string) (string, error) {
	compileCmd := fmt.Sprintf("g++ %s -o /tmp/cpp/program", filename)
	compileOutput, err := w.dockerContainer.ExecInContainerStdCopy(compileCmd)
	if compileOutput != "" {
		err = fmt.Errorf("compile error")
	}
	return compileOutput, err
}

func (w *Worker) execCpp(testcaseInput string) chan string {
	c := make(chan string)
	go func() {
		runCmd := fmt.Sprintf("echo '%s' > /tmp/input.txt && /tmp/cpp/program < /tmp/input.txt", testcaseInput)

		runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
		if err != nil {
			c <- err.Error()
			return
		}
		c <- runOutput
	}()
	return c
}

func (w *Worker) cleanUpCpp(filename string) {
	cleanupCmd := fmt.Sprintf("rm %s /tmp/cpp/program", filename)
	_, err := w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}

// func (w *Worker) runCppInContainer(job *Job, testcases []db.GetAllTestcasesRow) FinishedPayload {

// 	defer func() {
// 		cleanupCmd := fmt.Sprintf("rm %s /tmp/cpp/program", cppFileName)
// 		_, err = w.dockerContainer.ExecInContainer(cleanupCmd)
// 		if err != nil {
// 			log.Printf("Warning: failed to clean up files: %v", err)
// 		}
// 	}()

// 	// fetch all the input test cases
// 	// fetch all the required outputs
// 	var duration time.Duration = 0
// 	for _, val := range testcases {
// 		start := time.Now()
// 		runCmd := fmt.Sprintf("echo '%s' | /tmp/cpp/program", val.TestcaseInput.String)
// 		runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
// 		since := time.Since(start)
// 		duration += since
// 		if err != nil {

// 		}
// 		expectedLines := strings.Split(strings.TrimSpace(val.TestcaseOutput.String), "\n")
// 		actualLines := strings.Split(strings.TrimSpace(runOutput), "\n")
// 		for i := range expectedLines {
// 			expectedLines[i] = removeNonPrintableChars(strings.TrimSpace(expectedLines[i]))
// 		}
// 		for i := range actualLines {
// 			actualLines[i] = removeNonPrintableChars(strings.TrimSpace(actualLines[i]))
// 		}

// 		// if len(expectedLines) != len(actualLines) {
// 		// 	fmt.Println("Output mismatch: different number of lines")
// 		// 	return FinishedPayload{Message: "Wrong output", Where: "Output no of lines mismatch"}
// 		// }

// 		for i := range expectedLines {
// 			if expectedLines[i] != actualLines[i] {
// 				fmt.Printf("%#v\n", expectedLines[i])
// 				fmt.Printf("%#v\n", actualLines[i])
// 				return FinishedPayload{Message: "Wrong output", Where: "Testcase " + strconv.Itoa(int(val.TestcaseOrder.Int32+1)) + " Line no. " + strconv.Itoa(i+1) + ", expected " + expectedLines[i] + " getting  " + actualLines[i]}
// 			}
// 		}
// 	}

// 	return FinishedPayload{
// 		Message:   "Correct",
// 		TimeTaken: duration,
// 	}
// }
