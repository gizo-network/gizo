package chord

import (
	"errors"
	"fmt"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/kpango/glg"
)

const (
	MaxExecs = 10 // max number of jobs allowed in the chain
)

var (
	ErrJobsLenRange = errors.New("Number of jobs is more than allowed")
)

//Chord Jobs executed one after the other and the results passed to a callback
type Chord struct {
	jobs     []job.JobRequest
	bc       *core.BlockChain
	pq       *queue.JobPriorityQueue
	callback job.JobRequest
	result   job.JobRequest
	length   int
	status   string
}

//NewChord returns chord
func NewChord(j []job.JobRequest, callback job.JobRequest, bc *core.BlockChain, pq *queue.JobPriorityQueue) (*Chord, error) {
	//FIXME: count callback execs too
	length := len(callback.GetExec())
	for _, jr := range j {
		length += len(jr.GetExec())
	}

	if length > MaxExecs {
		return nil, ErrJobsLenRange
	}
	c := &Chord{
		jobs:     j,
		bc:       bc,
		pq:       pq,
		callback: callback,
		length:   length,
	}
	return c, nil
}

func (c Chord) GetCallback() job.JobRequest {
	return c.callback
}

func (c *Chord) setCallback(j job.JobRequest) {
	c.callback = j
}

//GetJobs returns jobs
func (c Chord) GetJobs() []job.JobRequest {
	return c.jobs
}

func (c *Chord) setJobs(j []job.JobRequest) {
	c.jobs = j
}

//GetStatus returns status
func (c Chord) GetStatus() string {
	return c.status
}

func (c *Chord) setStatus(s string) {
	c.status = s
}

func (c Chord) getLength() int {
	return c.length
}

func (c *Chord) setBC(bc *core.BlockChain) {
	c.bc = bc
}

func (c Chord) getPQ() *queue.JobPriorityQueue {
	return c.pq
}

func (c Chord) getBC() *core.BlockChain {
	return c.bc
}

func (c *Chord) setResults(res job.JobRequest) {
	c.result = res
}

//Result returns result
func (c Chord) Result() job.JobRequest {
	return c.result
}

//Dispatch executes the chord
func (c *Chord) Dispatch() {
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
				items = append(items, <-res) //! waits for it to finish before continuing
				c.getPQ().Push(*j, jr.GetExec()[i], res)
			}
		}
	}

	var callbackResults []queue.Item
	var callbackArgs []interface{} //holds result of execs
	for _, item := range items {
		callbackArgs = append(callbackArgs, item.GetExec().GetResult())
	}

	//! sets args as results or jobs
	for _, exec := range c.GetCallback().GetExec() {
		exec.SetArgs(callbackArgs)
	}

	j, err := c.getBC().FindJob(c.GetCallback().GetID())
	if err != nil {
		glg.Warn("Chain: Unable to find job - " + c.GetCallback().GetID())
		for _, exec := range c.GetCallback().GetExec() {
			exec.SetErr("Unable to find job - " + c.GetCallback().GetID())
		}
	} else {
		if len(c.GetCallback().GetExec()) > 1 {
			c.getPQ().Push(*j, c.GetCallback().GetExec()[0], res)
			for i := 1; i < len(c.GetCallback().GetExec()); i++ {
				callbackResults = append(callbackResults, <-res)
				c.getPQ().Push(*j, c.GetCallback().GetExec()[i], res)
			}
		} else {
			c.getPQ().Push(*j, c.GetCallback().GetExec()[0], res)
			callbackResults = append(callbackResults, <-res)
		}
	}

	fmt.Println("callbackargs", callbackArgs)
	fmt.Println("callbackResult", callbackResults[0].GetExec().GetResult())

	var callback job.JobRequest
	callback.SetID(c.GetCallback().GetID())
	for _, item := range callbackResults {
		callback.AppendExec(item.GetExec())
	}

	c.setResults(callback)
	c.setStatus(job.FINISHED)
}
