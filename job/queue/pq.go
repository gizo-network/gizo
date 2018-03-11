package queue

import (
	lane "github.com/Lobarr/lane"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue/qItem"
)

type JobPriorityQueue struct {
	pq *lane.PQueue
	bc *core.BlockChain
}

func (pq JobPriorityQueue) Push(j job.Job, exec *job.Exec, results chan<- qItem.Item, cancel chan struct{}) {
	pq.getPQ().Push(qItem.Item{
		Job: job.Job{
			ID:             j.GetID(),
			Hash:           j.GetHash(),
			Name:           j.GetName(),
			Task:           j.GetTask(),
			Signature:      j.GetSignature(),
			SubmissionTime: j.GetSubmissionTime(),
			Private:        j.GetPrivate(),
		},
		Exec:    exec,
		Results: results,
		Cancel:  cancel,
	}, exec.GetPriority())
}

func (pq JobPriorityQueue) Pop() qItem.Item {
	i, _ := pq.getPQ().Pop()
	return i.(qItem.Item)
}

func (pq JobPriorityQueue) Remove(hash []byte) {
	pq.pq.Remove(hash)
}

func (pq JobPriorityQueue) getPQ() *lane.PQueue {
	return pq.pq

}

func (pq JobPriorityQueue) watch() {
	for {
		if pq.getPQ().Empty() == false {
			//TODO: dispatch to next available worker node
			i := pq.Pop()
			if i.GetExec().GetStatus() == job.CANCELLED {
				i.ResultsChan() <- i
			} else {
				exec := i.Job.Execute(i.GetExec())
				i.SetExec(exec)
				i.ResultsChan() <- i
			}
		}
	}
}

func NewJobPriorityQueue() *JobPriorityQueue {
	pq := lane.NewPQueue(lane.MAXPQ)
	q := &JobPriorityQueue{
		pq: pq,
	}
	go q.watch()
	return q
}
