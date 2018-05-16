package job

import (
	"encoding/json"

	"github.com/kpango/glg"
)

type JobRequestSingle struct {
	ID   string
	Exec *Exec
}

func NewJobRequestSingle(id string, exec *Exec) *JobRequestSingle {
	return &JobRequestSingle{
		ID:   id,
		Exec: exec,
	}
}

func (jr *JobRequestSingle) SetID(id string) {
	jr.ID = id
}

func (jr JobRequestSingle) GetID() string {
	return jr.ID
}

func (jr JobRequestSingle) GetExec() *Exec {
	return jr.Exec
}

func (jr JobRequestSingle) Serialize() []byte {
	temp, err := json.Marshal(jr)
	if err != nil {
		glg.Error(err)
	}
	return temp
}

func DeserializeJRS(b []byte) (JobRequestSingle, error) {
	var temp JobRequestSingle
	err := json.Unmarshal(b, &temp)
	return temp, err
}
