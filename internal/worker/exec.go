package worker

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"ws-trial/db"
)

func (w Worker) Exec(job Job) FinishedPayload {
	// fetch input testcases, and output testcases
	query, pool, err := db.InitDb(w.ctx)
	defer pool.Close()
	if err != nil {
		fmt.Println("Error occured at db initialisation ", err.Error())
		return FinishedPayload{
			Message: "Error",
			Where:   "Some error occurred at the server",
		}
	}
	row, err := query.GetAllTestcases(w.ctx, int32(job.Qid))
	if err != nil {
		fmt.Print("Error in fetching the testcases ", err.Error())
		return FinishedPayload{
			Message: "Error",
			Where:   "Some error occured at the server",
		}
	}

	// compile and filecreation setup for all the languages
	var filename string
	var compileOut string

	switch job.Language {
	case "c++":
		filename = w.createCppFile(job.Code)
		compileOut, err = w.compileCpp(filename)
		defer w.cleanUpCpp(filename)
	case "java":
		filename = w.createJavaFile(job.Code)
		compileOut, err = w.compileJava(filename)
		defer w.cleanUpJava(filename)
	case "python":
		filename = w.createPythonFile(job.Code)
		defer w.cleanUpPython(filename)
	case "javascript":
		filename = w.createJavascriptFile(job.Code)
		defer w.cleanUpJavascript(filename)
	case "go":
		filename = w.createGoFile(job.Code)
		compileOut, err = w.compileGo(filename)
		defer w.cleanUpGo(filename)
	default:
	}
	if err != nil {
		return FinishedPayload{
			Message: "Error in compiling",
			Where:   compileOut,
		}
	}

	var outputChan chan string
	var outputString string
	var duration time.Duration

	for _, val := range row {
		start := time.Now()
		switch job.Language {
		case "c++":
			outputChan = w.execCpp(val.TestcaseInput.String)
		case "java":
			outputChan = w.execJava(val.TestcaseInput.String)
		case "python":
			outputChan = w.execPython(val.TestcaseInput.String, filename)
		case "javascript":
			outputChan = w.execJavascript(val.TestcaseInput.String, filename)
		case "go":
			outputChan = w.execGo(val.TestcaseInput.String)
		}

		select {
		case <-time.After(1 * time.Minute):
			w.dockerContainer.RestartContainer()
			return FinishedPayload{
				Message: "Your code took too long to execute",
			}
		case outputString = <-outputChan:
		}

		since := time.Since(start)
		duration += since

		expectedLines := strings.Split(strings.TrimSpace(val.TestcaseOutput.String), "\n")
		actualLines := strings.Split(strings.TrimSpace(outputString), "\n")
		for i := range expectedLines {
			expectedLines[i] = removeNonPrintableChars(strings.TrimSpace(expectedLines[i]))
		}
		for i := range actualLines {
			actualLines[i] = removeNonPrintableChars(strings.TrimSpace(actualLines[i]))
		}

		for i := range expectedLines {
			if expectedLines[i] != actualLines[i] {
				fmt.Printf("%#v\n", expectedLines[i])
				fmt.Printf("%#v\n", actualLines[i])
				return FinishedPayload{Message: "Wrong output", Where: "Testcase " + strconv.Itoa(int(val.TestcaseOrder.Int32+1)) + " Line no. " + strconv.Itoa(i+1) + ", expected " + expectedLines[i] + " getting  " + actualLines[i], TimeTaken: duration}
			}
		}
	}
	fmt.Print("It was correct!!!")
	return FinishedPayload{
		Message:   "Correct",
		TimeTaken: duration,
	}
}
