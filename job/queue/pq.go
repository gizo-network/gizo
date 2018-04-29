package queue

import (
	lane "github.com/Lobarr/lane"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue/qItem"
	"github.com/kpango/glg"
)

type JobPriorityQueue struct {
	pq *lane.PQueue
}

func (pq JobPriorityQueue) Push(j job.Job, exec *job.Exec, results chan<- qItem.Item, cancel chan struct{}) {
	temp := job.Job{
		ID:             j.GetID(),
		Hash:           j.GetHash(),
		Name:           j.GetName(),
		Task:           j.GetTask(),
		Signature:      j.GetSignature(),
		SubmissionTime: j.GetSubmissionTime(),
		Private:        j.GetPrivate(),
	}
	pq.GetPQ().Push(qItem.NewItem(temp, exec, results, cancel), exec.GetPriority())
	glg.Info("JobPriotityQueue: received job")
}

func (pq JobPriorityQueue) PushItem(i qItem.Item, piority int) {
	pq.GetPQ().Push(i, piority)
	glg.Info("JobPriotityQueue: received job")

}

func (pq JobPriorityQueue) Pop() qItem.Item {
	i, _ := pq.GetPQ().Pop()
	return i.(qItem.Item)
}

func (pq JobPriorityQueue) Remove(hash []byte) {
	pq.pq.RemoveHash(hash)
}

func (pq JobPriorityQueue) GetPQ() *lane.PQueue {
	return pq.pq

}

// func (pq JobPriorityQueue) watch() {
// 	for {
// 		if pq.getPQ().Empty() == false {
// 			//TODO: dispatch to next available worker node
// 			i := pq.Pop()
// 			if i.GetExec().GetStatus() == job.CANCELLED {
// 				i.ResultsChan() <- i
// 			} else {
// 				exec := i.Job.Execute(i.GetExec())
// 				i.SetExec(exec)
// 				i.ResultsChan() <- i
// 			}
// 		}
// 	}
// }

func NewJobPriorityQueue() *JobPriorityQueue {
	pq := lane.NewPQueue(lane.MAXPQ)
	q := &JobPriorityQueue{
		pq: pq,
	}
	// go q.watch()
	return q
}
