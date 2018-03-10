package chain

import (
	"sync"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/gizo-network/gizo/job/queue/qItem"
)

//Chain - Jobs executed one after the other
type Chain struct {
	jobs   []job.JobRequestMultiple
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	result []job.JobRequestMultiple
	length int
	status string
	cancel chan struct{}
}

//NewChain returns chain
func NewChain(j []job.JobRequestMultiple, bc *core.BlockChain, pq *queue.JobPriorityQueue) (*Chain, error) {
	length := 0
	for _, jr := range j {
		length += len(jr.GetExec())
	}
	if length > job.MaxExecs {
		return nil, job.ErrJobsLenRange
	}
	c := &Chain{
		jobs:   j,
		bc:     bc,
		pq:     pq,
		length: length,
		cancel: make(chan struct{}),
	}
	return c, nil
}

func (c *Chain) Cancel() {
	c.cancel <- struct{}{}
}

func (c Chain) GetCancelChan() chan struct{} {
	return c.cancel
}

//GetJobs returns jobs
func (c Chain) GetJobs() []job.JobRequestMultiple {
	return c.jobs
}

func (c *Chain) setJobs(j []job.JobRequestMultiple) {
	c.jobs = j
}

//GetStatus returns status
func (c Chain) GetStatus() string {
	return c.status
}

func (c *Chain) setStatus(s string) {
	c.status = s
}

func (c Chain) getLength() int {
	return c.length
}

func (c *Chain) setBC(bc *core.BlockChain) {
	c.bc = bc
}

func (c Chain) getPQ() *queue.JobPriorityQueue {
	return c.pq
}

func (c Chain) getBC() *core.BlockChain {
	return c.bc
}

func (c *Chain) setResults(res []job.JobRequestMultiple) {
	c.result = res
}

//Result returns result
func (c Chain) Result() []job.JobRequestMultiple {
	return c.result
}

//Dispatch executes the chain
func (c *Chain) Dispatch() {
	c.setStatus(job.RUNNING)
	var results []qItem.Item // used to hold results
	res := make(chan qItem.Item)
	cancelled := false
	closeCancel := make(chan struct{})
	var wg sync.WaitGroup
	//! watch cancel channel
	wg.Add(1)
	go func() {
		select {
		case <-c.cancel:
			cancelled = true
			glg.Warn("Chain: Cancelling jobs")
			for _, jr := range c.GetJobs() {
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
	var jobIDs []string
	for _, jr := range c.GetJobs() {
		c.setStatus("Queueing execs of job - " + jr.GetID())
		jobIDs = append(jobIDs, jr.GetID())
		j, err := c.getBC().FindJob(jr.GetID())
		if err != nil {
			glg.Warn("Chain: Unable to find job - " + jr.GetID())
			for _, exec := range jr.GetExec() {
				exec.SetErr("Unable to find job - " + jr.GetID())
			}
		} else {
			for i := 0; i < len(jr.GetExec()); i++ {
				if cancelled == true {
					results = append(results, qItem.Item{
						Job: job.Job{
							ID:             j.GetID(),
							Hash:           j.GetHash(),
							Name:           j.GetName(),
							Task:           j.GetTask(),
							Signature:      j.GetSignature(),
							SubmissionTime: j.GetSubmissionTime(),
							Private:        j.GetPrivate(),
						},
						Exec:    jr.GetExec()[i],
						Results: res,
						Cancel:  c.GetCancelChan(),
					})
				} else {
					c.getPQ().Push(*j, jr.GetExec()[i], res, c.GetCancelChan()) //? queues first job
					results = append(results, <-res)
				}
			}
		}
	}
	close(res)

	var grouped []job.JobRequestMultiple
	for _, jID := range jobIDs {
		var req job.JobRequestMultiple
		req.SetID(jID)
		for _, item := range results {
			if item.GetID() == jID {
				req.AppendExec(item.GetExec())
			}
		}
		grouped = append(grouped, req)
	}

	if cancelled == false {
		closeCancel <- struct{}{}
		c.setStatus(job.FINISHED)
	} else {
		c.setStatus(job.CANCELLED)
	}
	wg.Wait()
	c.setResults(grouped)
}
