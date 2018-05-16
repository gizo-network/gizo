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

func (pq WorkerPriorityQueue) Remove(s *melody.Session) {
	pq.getPQ().RemoveSession(s)
}

func (pq WorkerPriorityQueue) getPQ() *lane.PQueue {
	return pq.pq
}

func (pq WorkerPriorityQueue) Len() int {
	return pq.getPQ().Size()
}

func NewWorkerPriorityQueue() *WorkerPriorityQueue {
	pq := lane.NewPQueue(lane.MINPQ)
	q := &WorkerPriorityQueue{
		pq: pq,
	}
	return q
}
