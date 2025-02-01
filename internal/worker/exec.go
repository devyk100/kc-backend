package worker

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"
	"ws-trial/db"
)

func FilterTestcases(testcaseOut string, actualOut string) ([]string, []string) {
	expectedLines := strings.Split(strings.TrimSpace(testcaseOut), "\n")
	actualLines := strings.Split(strings.TrimSpace(actualOut), "\n")
	filteredExpected := []string{}
	for _, line := range expectedLines {
		cleaned := removeNonPrintableChars(strings.TrimSpace(line))
		if cleaned != "" {
			filteredExpected = append(filteredExpected, cleaned)
		}
	}

	filteredActual := []string{}
	for _, line := range actualLines {
		cleaned := removeNonPrintableChars(strings.TrimSpace(line))
		if cleaned != "" {
			filteredActual = append(filteredActual, cleaned)
		}
	}

	return filteredExpected, filteredActual
}

func (w Worker) Exec(job Job, query *db.Queries) FinishedPayload {
	job.Code = html.UnescapeString(job.Code)
	// fetch input testcases, and output testcases
	row, err := query.GetAllTestcases(w.ctx, int32(job.Qid))
	if err != nil {
		fmt.Print("Error in fetching the testcases ", err.Error(), "\n\n")
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
		fmt.Println("Actual output", outputString)
		outputString = removeNonPrintableChars(outputString)
		expectedLines, actualLines := FilterTestcases(val.TestcaseOutput.String, outputString)
		i := 0
		for ; i < len(expectedLines) && i < len(actualLines); i++ {
			if expectedLines[i] != actualLines[i] {
				fmt.Printf("%#v\n", expectedLines[i])
				fmt.Printf("%#v\n", actualLines[i])
				return FinishedPayload{Message: "Wrong output", Where: "Testcase " + strconv.Itoa(int(val.TestcaseOrder.Int32+1)) + " Line no. " + strconv.Itoa(i+1) + ", expected " + expectedLines[i] + " getting  " + actualLines[i], TimeTaken: int32(duration.Milliseconds())}
			}
		}
		if i < len(expectedLines) {
			return FinishedPayload{Message: "Wrong output", Where: "Testcase no of output is not same" + val.TestcaseOutput.String + " your out " + outputString, TimeTaken: int32(duration.Milliseconds())}
		}
	}
	fmt.Print("It was correct!!!")
	return FinishedPayload{
		Message:   "Correct",
		TimeTaken: int32(duration.Milliseconds()),
	}
}
