package job

type JobRequestSingle struct {
	id   string
	exec *Exec
}

func NewJobRequestSingle(id string, exec *Exec) *JobRequestSingle {
	return &JobRequestSingle{
		id:   id,
		exec: exec,
	}
}

func (jr *JobRequestSingle) SetID(id string) {
	jr.id = id
}

func (jr JobRequestSingle) GetID() string {
	return jr.id
}

func (jr JobRequestSingle) GetExec() *Exec {
	return jr.exec
}
