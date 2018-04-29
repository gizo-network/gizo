package p2p

import "github.com/gizo-network/gizo/job/queue/qItem"

type WorkerInfo struct {
	pub string
	job *qItem.Item
}
