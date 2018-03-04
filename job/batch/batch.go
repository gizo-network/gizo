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

type Batch struct {
	jobs   []job.JobRequest
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	result []job.JobRequest
	length int
	status string
}

func NewBatch(j []job.JobRequest, bc *core.BlockChain, pq *queue.JobPriorityQueue) (*Batch, error) {
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

func (b Batch) GetJobs() []job.JobRequest {
	return b.jobs
}

func (b *Batch) SetJobs(j []job.JobRequest) {
	b.jobs = j
}

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

func (b *Batch) setResults(res []job.JobRequest) {
	b.result = res
}

func (b Batch) Result() []job.JobRequest {
	return b.result
}

func (b *Batch) Dispatch() {
	//! should be run as a go routine because it blocks till all jobs are complete
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

	var grouped []job.JobRequest
	for _, jID := range jobIDs {
		var req job.JobRequest
		req.SetID(jID)
		for item := range results {
			if item.GetID() == jID {
				req.AppendExec(item.GetExec())
			}
		}
		grouped = append(grouped, req)
	}
	b.setResults(grouped)
	b.setStatus(job.FINISHED)
}
