package qItem

import "github.com/gizo-network/gizo/job"

type Item struct {
	Job     job.Job
	Exec    *job.Exec
	Results chan<- Item
	Cancel  chan<- struct{}
}

func (i Item) GetCancel() chan<- struct{} {
	return i.Cancel
}

//sets exec
func (i *Item) SetExec(ex *job.Exec) {
	i.Exec = ex
}

//GetExec returns exec
func (i Item) GetExec() *job.Exec {
	return i.Exec
}

//GetID returns id
func (i Item) GetID() string {
	return i.Job.GetID()
}

//GetJob returns job
func (i Item) GetJob() job.Job {
	return i.Job
}

//ResultsChan return result chan
func (i Item) ResultsChan() chan<- Item {
	return i.Results
}
