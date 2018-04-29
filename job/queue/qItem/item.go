package qItem

import (
	"encoding/json"

	"github.com/gizo-network/gizo/job"
	"github.com/kpango/glg"
)

type Item struct {
	Job     job.Job   `json:"job"`
	Exec    *job.Exec `json:"exec"`
	results chan<- Item
	cancel  chan<- struct{}
}

func NewItem(j job.Job, exec *job.Exec, results chan<- Item, cancel chan<- struct{}) Item {
	return Item{
		Job:     j,
		Exec:    exec,
		results: results,
		cancel:  cancel,
	}
}

func (i Item) GetCancel() chan<- struct{} {
	return i.cancel
}

//sets exec
func (i *Item) SetExec(ex *job.Exec) {
	i.Exec = ex
}

//GetExec returns exec
func (i Item) GetExec() *job.Exec {
	return i.Exec
}

//GetID returns id
func (i Item) GetID() string {
	return i.Job.GetID()
}

//GetJob returns job
func (i Item) GetJob() job.Job {
	return i.Job
}

//ResultsChan return result chan
func (i Item) ResultsChan() chan<- Item {
	return i.results
}

func (i Item) Serialize() []byte {
	bytes, err := json.Marshal(i)
	if err != nil {
		glg.Fatal(err)
	}
	return bytes
}

func DeserializeItem(b []byte) Item {
	var temp Item
	err := json.Unmarshal(b, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}
