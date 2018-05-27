package solo

import (
	"strconv"
	"sync"
	"time"

	"github.com/gizo-network/gizo/cache"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/gizo-network/gizo/job/queue/qItem"
	"github.com/kpango/glg"
)

//Solo - Jobs executed one after the other
type Solo struct {
	jr     job.JobRequestSingle
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	jc     *cache.JobCache
	result job.JobRequestSingle
	status string
	cancel chan struct{}
}

func NewSolo(jr job.JobRequestSingle, bc *core.BlockChain, pq *queue.JobPriorityQueue, jc *cache.JobCache) *Solo {
	return &Solo{
		jr: jr,
		bc: bc,
		pq: pq,
		jc: jc,
	}
}

func (s *Solo) Cancel() {
	s.cancel <- struct{}{}
}

func (s Solo) GetCancelChan() chan struct{} {
	return s.cancel
}

func (s Solo) GetJob() job.JobRequestSingle {
	return s.jr
}

func (s Solo) GetStatus() string {
	return s.status
}

func (s *Solo) setStatus(status string) {
	s.status = status
}

func (s *Solo) setBC(bc *core.BlockChain) {
	s.bc = bc
}

func (s Solo) getPQ() *queue.JobPriorityQueue {
	return s.pq
}

func (s Solo) getJC() *cache.JobCache {
	return s.jc
}

func (s Solo) getBC() *core.BlockChain {
	return s.bc
}

func (s *Solo) setResult(res job.JobRequestSingle) {
	s.result = res
}

//Result returns result
func (s Solo) Result() job.JobRequestSingle {
	return s.result
}

//Dispatch executes the solo
func (s *Solo) Dispatch() {
	s.setStatus(job.RUNNING)
	var result qItem.Item
	res := make(chan qItem.Item)
	cancelled := false
	closeCancel := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case <-s.cancel:
			cancelled = true
			glg.Warn("Solo: Cancelling job")
			if s.GetJob().GetExec().GetStatus() == job.RUNNING || s.GetJob().GetExec().GetStatus() == job.RETRYING {
				s.GetJob().GetExec().Cancel()
			}
			if s.GetJob().GetExec().GetResult() == nil {
				s.GetJob().GetExec().SetStatus(job.CANCELLED)
			}
			break
		case <-closeCancel:
			break
		}
		wg.Done()
	}()
	s.setStatus("Queueing execs of job - " + s.GetJob().GetID())
	var j *job.Job
	var err error
	j, err = s.getJC().Get(s.GetJob().GetID()) //check job cache first
	if j == nil {
		j, err = s.getBC().FindJob(s.GetJob().GetID())
	}
	if err != nil {
		glg.Warn("Batch: Unable to find job - " + s.GetJob().GetID())
		s.GetJob().GetExec().SetErr("Batch: Unable to find job - " + s.GetJob().GetID())
	} else {
		if cancelled == true {
			result = qItem.NewItem(job.Job{
				ID:             j.GetID(),
				Hash:           j.GetHash(),
				Name:           j.GetName(),
				Task:           j.GetTask(),
				Signature:      j.GetSignature(),
				SubmissionTime: j.GetSubmissionTime(),
				Private:        j.GetPrivate(),
			}, s.GetJob().GetExec(), res, s.GetCancelChan())
		} else {
			if s.GetJob().GetExec().GetExecutionTime() != 0 {
				glg.Warn("Chord: Queuing in " + strconv.FormatFloat(time.Unix(s.GetJob().GetExec().GetExecutionTime(), 0).Sub(time.Now()).Seconds(), 'f', -1, 64) + " nanoseconds")
				time.Sleep(time.Nanosecond * time.Duration(time.Unix(s.GetJob().GetExec().GetExecutionTime(), 0).Sub(time.Now()).Nanoseconds()))
			}
			s.getPQ().Push(*j, s.GetJob().GetExec(), res, s.GetCancelChan())
			result = <-res
		}
	}
	close(res)

	if cancelled == false {
		closeCancel <- struct{}{}
		s.setStatus(job.FINISHED)
	} else {
		s.setStatus(job.CANCELLED)
	}
	wg.Wait()
	s.setResult(*job.NewJobRequestSingle(result.GetID(), result.GetExec()))
}
