package chain

import (
	"errors"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
)

const (
	MaxExecs = 10 // max number of jobs allowed in the chain
)

var (
	ErrJobsLenRange = errors.New("Number of jobs is more than allowed")
)

//Chain - Jobs executed one after the other
type Chain struct {
	jobs   []job.JobRequest
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	result []job.JobRequest
	length int
	status string
}

//NewChain returns chain
func NewChain(j []job.JobRequest, bc *core.BlockChain, pq *queue.JobPriorityQueue) (*Chain, error) {
	length := 0
	for _, jr := range j {
		length += len(jr.GetExec())
	}
	if length > MaxExecs {
		return nil, ErrJobsLenRange
	}
	c := &Chain{
		jobs:   j,
		bc:     bc,
		pq:     pq,
		length: length,
	}
	return c, nil
}

//GetJobs returns jobs
func (c Chain) GetJobs() []job.JobRequest {
	return c.jobs
}

func (c *Chain) setJobs(j []job.JobRequest) {
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

func (c *Chain) setResults(res []job.JobRequest) {
	c.result = res
}

//Result returns result
func (c Chain) Result() []job.JobRequest {
	return c.result
}

//Dispatch executes the chain
func (c *Chain) Dispatch() {
	c.setStatus(job.RUNNING)
	var items []queue.Item // used to hold results
	res := make(chan queue.Item)
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
			c.getPQ().Push(*j, jr.GetExec()[0], res) //? queues first job
			for i := 1; i < len(jr.GetExec()); i++ {
				items = append(items, <-res)
				c.getPQ().Push(*j, jr.GetExec()[i], res)
			}
		}
	}

	var grouped []job.JobRequest
	for _, jID := range jobIDs {
		var req job.JobRequest
		req.SetID(jID)
		for _, item := range items {
			if item.GetID() == jID {
				req.AppendExec(item.GetExec())
			}
		}
		grouped = append(grouped, req)
	}
	c.setResults(grouped)
	c.setStatus(job.FINISHED)
}
