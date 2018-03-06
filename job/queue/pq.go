package queue

import (
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	lane "gopkg.in/oleiade/lane.v1"
)

type JobPriorityQueue struct {
	pq *lane.PQueue
	bc *core.BlockChain
}

func (pq JobPriorityQueue) Push(j job.Job, exec *job.Exec, results chan<- Item) {
	pq.getPQ().Push(Item{
		job: job.Job{
			ID:             j.GetID(),
			Hash:           j.GetHash(),
			Name:           j.GetName(),
			Task:           j.GetTask(),
			Signature:      j.GetSignature(),
			SubmissionTime: j.GetSubmissionTime(),
			Private:        j.GetPrivate(),
		},
		exec:    exec,
		results: results,
	}, exec.GetPriority())
}

func (pq JobPriorityQueue) Pop() Item {
	item, _ := pq.getPQ().Pop()
	return item.(Item)
}

func (pq JobPriorityQueue) getPQ() *lane.PQueue {
	return pq.pq
}

func (pq JobPriorityQueue) watch() {
	go func() {
		for {
			if pq.getPQ().Empty() == false {
				//TODO: dispatch to next available worker node
				item := pq.Pop()
				exec := item.job.Execute(item.GetExec())
				item.setExec(exec)
				item.ResultsChan() <- item
			}
		}
	}()
}

func NewJobPriorityQueue() *JobPriorityQueue {
	pq := lane.NewPQueue(lane.MAXPQ)
	q := &JobPriorityQueue{
		pq: pq,
	}
	q.watch()
	return q
}
