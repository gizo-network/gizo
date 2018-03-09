package batch

import (
	"errors"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/kpango/glg"
)

const (
	MaxExecs = 10 // max number of jobs allowed in the Batch
)

var (
	ErrJobsLenRange = errors.New("Number of jobs is more than allowed")
)

//Batch - Jobs executed in parralele
type Batch struct {
	jobs   []job.JobRequestMultiple
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	result []job.JobRequestMultiple
	length int
	status string
}

//NewBatch returns batch
func NewBatch(j []job.JobRequestMultiple, bc *core.BlockChain, pq *queue.JobPriorityQueue) (*Batch, error) {
	length := 0
	for _, jr := range j {
		length += len(jr.GetExec())
	}
	if length > MaxExecs {
		return nil, ErrJobsLenRange
	}
	b := &Batch{
		jobs:   j,
		bc:     bc,
		pq:     pq,
		length: length,
	}

	return b, nil
}

//GetJobs return jobs
func (b Batch) GetJobs() []job.JobRequestMultiple {
	return b.jobs
}

func (b *Batch) setJobs(j []job.JobRequestMultiple) {
	b.jobs = j
}

//GetStatus returns status
func (b Batch) GetStatus() string {
	return b.status
}

func (b *Batch) setStatus(s string) {
	b.status = s
}

func (b *Batch) setBC(bc *core.BlockChain) {
	b.bc = bc
}

func (b Batch) getBC() *core.BlockChain {
	return b.bc
}

func (b Batch) getPQ() *queue.JobPriorityQueue {
	return b.pq
}

func (b Batch) getLength() int {
	return b.length
}

func (b *Batch) setResults(res []job.JobRequestMultiple) {
	b.result = res
}

//Result returns result
func (b Batch) Result() []job.JobRequestMultiple {
	return b.result
}

//Dispatch executes the batch
func (b *Batch) Dispatch() {
	//! should be run in a go routine because it blocks till all jobs are complete
	var items []queue.Item
	b.setStatus(job.RUNNING)
	results := make(chan queue.Item, b.getLength())
	var jobIDs []string
	for _, jr := range b.GetJobs() {
		b.setStatus("Queueing execs of job - " + jr.GetID())
		jobIDs = append(jobIDs, jr.GetID())
		j, err := b.getBC().FindJob(jr.GetID())
		if err != nil {
			glg.Warn("Batch: Unable to find job - " + jr.GetID())
			for _, exec := range jr.GetExec() {
				exec.SetErr("Unable to find job - " + jr.GetID())
			}
		} else {
			for _, exec := range jr.GetExec() {
				b.getPQ().Push(*j, exec, results)
			}
		}
	}

	//! wait for all jobs to be done
	for {
		if len(results) == cap(results) {
			close(results)
			break
		}
	}

	for item := range results {
		items = append(items, item)
	}

	var grouped []job.JobRequestMultiple
	for _, jID := range jobIDs {
		var req job.JobRequestMultiple
		req.SetID(jID)
		for _, item := range items {
			if item.GetID() == jID {
				req.AppendExec(item.GetExec())
			}
		}
		grouped = append(grouped, req)
	}
	b.setResults(grouped)
	b.setStatus(job.FINISHED)
}
