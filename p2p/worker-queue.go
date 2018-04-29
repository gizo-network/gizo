package p2p

import (
	"github.com/Lobarr/lane"
	"gopkg.in/olahol/melody.v1"
)

type WorkerPriorityQueue struct {
	pq *lane.PQueue
}

func (pq WorkerPriorityQueue) Push(s *melody.Session) {

}
