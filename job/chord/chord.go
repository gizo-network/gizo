package chord

import (
	"errors"

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
	jobs     []job.JobRequestMultiple
	bc       *core.BlockChain
	pq       *queue.JobPriorityQueue
	callback job.JobRequestMultiple
	result   job.JobRequestMultiple
	length   int
	status   string
}

//NewChord returns chord
func NewChord(j []job.JobRequestMultiple, callback job.JobRequestMultiple, bc *core.BlockChain, pq *queue.JobPriorityQueue) (*Chord, error) {
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

func (c Chord) GetCallback() job.JobRequestMultiple {
	return c.callback
}

func (c *Chord) setCallback(j job.JobRequestMultiple) {
	c.callback = j
}

//GetJobs returns jobs
func (c Chord) GetJobs() []job.JobRequestMultiple {
	return c.jobs
}

func (c *Chord) setJobs(j []job.JobRequestMultiple) {
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

func (c *Chord) setResults(res job.JobRequestMultiple) {
	c.result = res
}

//Result returns result
func (c Chord) Result() job.JobRequestMultiple {
	return c.result
}

//Dispatch executes the chord
func (c *Chord) Dispatch() {
	c.setStatus(job.RUNNING)
	var items []queue.Item // used to hold results
	resChan := make(chan queue.Item, 1)
	for _, jr := range c.GetJobs() {
		c.setStatus("Queueing execs of job - " + jr.GetID())
		j, err := c.getBC().FindJob(jr.GetID())
		if err != nil {
			glg.Warn("Chord: Unable to find job - " + jr.GetID())
			for _, exec := range jr.GetExec() {
				exec.SetErr("Unable to find job - " + jr.GetID())
			}
		} else {
			for i := 0; i < len(jr.GetExec()); i++ {
				c.getPQ().Push(*j, jr.GetExec()[i], resChan)
				items = append(items, <-resChan) //! waits for it to finish before continuing
			}
		}
	}
	close(resChan)

	var callbackResults []queue.Item
	var callbackArgs []interface{}           //holds result of execs
	callbackChan := make(chan queue.Item, 1) //causes program to pause
	for _, item := range items {
		callbackArgs = append(callbackArgs, item.GetExec().GetResult())
	}

	//! sets args as results or jobs
	for _, exec := range c.GetCallback().GetExec() {
		exec.SetArgs(callbackArgs)
	}

	cj, err := c.getBC().FindJob(c.GetCallback().GetID())
	if err != nil {
		glg.Warn("Chord: Unable to find job - " + c.GetCallback().GetID())
		for _, exec := range c.GetCallback().GetExec() {
			exec.SetErr("Unable to find job - " + c.GetCallback().GetID())
		}
	} else {
		if len(c.GetCallback().GetExec()) > 1 {
			c.getPQ().Push(*cj, c.GetCallback().GetExec()[0], callbackChan)
			for i := 1; i < len(c.GetCallback().GetExec()); i++ {
				callbackResults = append(callbackResults, <-callbackChan)
				c.getPQ().Push(*cj, c.GetCallback().GetExec()[i], callbackChan)
			}
		} else {
			c.getPQ().Push(*cj, c.GetCallback().GetExec()[0], callbackChan)
			callbackResults = append(callbackResults, <-callbackChan)
		}
	}

	close(callbackChan)

	var callback job.JobRequestMultiple
	callback.SetID(c.GetCallback().GetID())
	for _, item := range callbackResults {
		callback.AppendExec(item.GetExec())
	}

	c.setResults(callback)
	c.setStatus(job.FINISHED)
}
