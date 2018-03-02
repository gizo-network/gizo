package chain

import (
	"errors"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
)

const (
	MaxLen = 10
)

var (
	ErrJobsLenRange = errors.New("Number of jobs is more than allowed")
)

type Chain struct {
	jobs   []job.JobRequest
	bc     *core.BlockChain
	status string
}

func NewChain(j []job.JobRequest, bc *core.BlockChain) (*Chain, error) {
	if len(j) > MaxLen {
		return nil, ErrJobsLenRange
	}
	c := &Chain{}
	c.SetJobs(j)
	c.setBC(bc)
	return c, nil
}

func (c Chain) GetJobs() []job.JobRequest {
	return c.jobs
}

func (c *Chain) SetJobs(j []job.JobRequest) {
	c.jobs = j
}

func (c Chain) GetStatus() string {
	return c.status
}

func (c *Chain) setStatus(s string) {
	c.status = s
}

func (c *Chain) setBC(bc *core.BlockChain) {
	c.bc = bc
}

func (c Chain) getBC() *core.BlockChain {
	return c.bc
}

func (c Chain) Dispatch() {
	c.setStatus(job.RUNNING)
	for _, jr := range c.GetJobs() {
		c.setStatus("Running job - " + jr.GetID())
		job, err := c.getBC().FindJob(jr.GetID())
		if err != nil {
			glg.Warn("Chain: Unable to find job - " + jr.GetID())
			for _, exec := range jr.GetExec() {
				exec.SetErr("Unable to find job - " + jr.GetID())
			}
		} else {
			for _, exec := range jr.GetExec() {
				job.Execute(exec) //! add to queue
			}
		}
	}
	c.setStatus(job.FINISHED)
}

func (c Chain) Result() [][]*job.Exec {
	var temp [][]*job.Exec
	for _, jr := range c.GetJobs() {
		temp = append(temp, jr.GetExec())
	}
	return temp
}
