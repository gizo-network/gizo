package queue

import "github.com/gizo-network/gizo/job"

type Item struct {
	job     job.Job
	exec    *job.Exec
	results chan<- Item
}

func (i *Item) SetExec(ex *job.Exec) {
	i.exec = ex
}

func (i Item) GetExec() *job.Exec {
	return i.exec
}

func (i Item) GetID() string {
	return i.job.GetID()
}

func (i Item) GetJob() job.Job {
	return i.job
}

func (i Item) ResultsChan() chan<- Item {
	return i.results
}
