package queue

import "github.com/gizo-network/gizo/job"

type Item struct {
	job     job.Job
	exec    *job.Exec
	results chan<- Item
}

//sets exec
func (i *Item) setExec(ex *job.Exec) {
	i.exec = ex
}

//GetExec returns exec
func (i Item) GetExec() *job.Exec {
	return i.exec
}

//GetID returns id
func (i Item) GetID() string {
	return i.job.GetID()
}

//GetJob returns job
func (i Item) GetJob() job.Job {
	return i.job
}

//ResultsChan return result chan
func (i Item) ResultsChan() chan<- Item {
	return i.results
}
