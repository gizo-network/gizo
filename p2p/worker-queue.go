package p2p

import (
	"github.com/Lobarr/lane"
	"gopkg.in/olahol/melody.v1"
)

type WorkerPriorityQueue struct {
	pq *lane.PQueue
}

func (pq WorkerPriorityQueue) Push(s *melody.Session, priority int) {
	pq.getPQ().Push(s, priority)
}

func (pq WorkerPriorityQueue) Pop() *melody.Session {
	i, _ := pq.getPQ().Pop()
	return i.(*melody.Session)
}

func (pq WorkerPriorityQueue) Remove(hash []byte) {
	pq.getPQ().RemoveHash(hash)
}

func (pq WorkerPriorityQueue) getPQ() *lane.PQueue {
	return pq.pq

}

// func (pq WorkerPriorityQueue) watch() {
// 	for {
// 		if pq.getPQ().Empty() == false {
// 			//TODO: dispatch to next available worker node
// 			// i := pq.Pop()
// 		}
// 	}
// }

func NewWorkerPriorityQueue() *WorkerPriorityQueue {
	pq := lane.NewPQueue(lane.MINPQ)
	q := &WorkerPriorityQueue{
		pq: pq,
	}
	return q
}
