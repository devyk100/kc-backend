package worker

func (w Worker) Exec(job Job) {
	switch job.Language {
	case "c++":
		w.runCppInContainer(&job)
	case "java":
	case "python":
	case "javascript":
	case "go":
		w.runGoInContainer(&job)
	default:
	}
}
