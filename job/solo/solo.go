package solo

import (
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/kpango/glg"
)

//Solo - Jobs executed one after the other
type Solo struct {
	jr     job.JobRequestSingle
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	result job.JobRequestSingle
	status string
}

func NewSolo(jr job.JobRequestSingle, bc *core.BlockChain, pq *queue.JobPriorityQueue) *Solo {
	return &Solo{
		jr: jr,
		bc: bc,
		pq: pq,
	}
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

func (s *Solo) Dispatch() {
	s.setStatus(job.RUNNING)
	var result queue.Item
	res := make(chan queue.Item)
	s.setStatus("Queueing execs of job - " + s.GetJob().GetID())
	j, err := s.getBC().FindJob(s.GetJob().GetID())
	if err != nil {
		glg.Warn("Batch: Unable to find job - " + s.GetJob().GetID())
		s.GetJob().GetExec().SetErr("Batch: Unable to find job - " + s.GetJob().GetID())
	} else {
		s.getPQ().Push(*j, s.GetJob().GetExec(), res)
		result = <-res
	}
	close(res)

	s.setResult(*job.NewJobRequestSingle(result.GetID(), result.GetExec()))
	s.setStatus(job.FINISHED)
}
