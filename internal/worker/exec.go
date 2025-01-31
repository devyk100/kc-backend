package worker

import "fmt"

func (w Worker) Exec(job Job) {
	switch job.Language {
	case "c++":
		w.runCppInContainer(&job)
		break
	case "java":
		w.runJavaInContainer(&job)
		break
	case "python":
		w.runPythonInContainer(&job)
		break
	case "javascript":
		w.runJavaScriptInContainer(&job)
		break
	case "go":
		w.runGoInContainer(&job)
		break
	default:
	}
	fmt.Println("Tried executing the code")

}
