package batch

import (
	"sync"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/gizo-network/gizo/job/queue/qItem"
	"github.com/kpango/glg"
)

//Batch - Jobs executed in parralele
type Batch struct {
	jobs   []job.JobRequestMultiple
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	result []job.JobRequestMultiple
	length int
	status string
	cancel chan struct{}
}

//NewBatch returns batch
func NewBatch(j []job.JobRequestMultiple, bc *core.BlockChain, pq *queue.JobPriorityQueue) (*Batch, error) {
	length := 0
	for _, jr := range j {
		length += len(jr.GetExec())
	}
	if length > job.MaxExecs {
		return nil, job.ErrJobsLenRange
	}
	b := &Batch{
		jobs:   j,
		bc:     bc,
		pq:     pq,
		length: length,
		cancel: make(chan struct{}),
	}

	return b, nil
}

func (b *Batch) Cancel() {
	b.cancel <- struct{}{}
}

func (b Batch) GetCancelChan() chan struct{} {
	return b.cancel
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
	var items []qItem.Item
	b.setStatus(job.RUNNING)
	cancelled := false
	closeCancel := make(chan struct{})
	var wg sync.WaitGroup
	//! watch cancel channel
	wg.Add(1)
	go func() {
		select {
		case <-b.cancel:
			cancelled = true
			glg.Warn("Batch: Cancelling jobs")
			for _, jr := range b.GetJobs() {
				for _, exec := range jr.GetExec() {
					if exec.GetStatus() == job.RUNNING || exec.GetStatus() == job.RETRYING {
						exec.Cancel()
					}
					if exec.GetResult() == nil {
						exec.SetStatus(job.CANCELLED)
					}
				}
			}
			break
		case <-closeCancel:
			break
		}
		wg.Done()
	}()

	results := make(chan qItem.Item, b.getLength())
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
				b.getPQ().Push(*j, exec, results, b.GetCancelChan())
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
	if cancelled == false {
		closeCancel <- struct{}{}
		b.setStatus(job.FINISHED)
	} else {
		b.setStatus(job.CANCELLED)
	}
	wg.Wait()
	b.setResults(grouped)
}
