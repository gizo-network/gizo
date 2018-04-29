package p2p

import "github.com/gizo-network/gizo/job/queue/qItem"

type WorkerInfo struct {
	pub string
	job *qItem.Item
}

func NewWorkerInfo(pub string) *WorkerInfo {
	return &WorkerInfo{pub: pub}
}

func (w WorkerInfo) GetPub() string {
	return w.pub
}

func (w *WorkerInfo) SetPub(pub string) {
	w.pub = pub
}

func (w WorkerInfo) GetJob() *qItem.Item {
	return w.job
}

func (w *WorkerInfo) SetJob(j *qItem.Item) {
	w.job = j
}

func (w *WorkerInfo) Assign(j *qItem.Item) {
	w.job = j
}

func (w *WorkerInfo) Busy() bool {
	return w.GetJob() == nil
}
