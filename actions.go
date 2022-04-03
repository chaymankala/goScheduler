package main

func doActionType(atype int, job *Job) {
	if atype == 1 {
		sendAPI(job)
	}
}

func sendAPI(job *Job) {
	_, err := postHTTP("https://google.com", makeParams(job.Meta[0].Value, job.JobName))
	if err == nil {
		removeJob(job.Id)
	}
}
