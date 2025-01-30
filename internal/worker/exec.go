package worker

func (w Worker) Exec(job Job) {
	switch job.Language {
	case "c++":
	case "java":
	case "python":
	case "javascript":
	case "go":
	default:
	}
}
