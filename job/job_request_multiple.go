package job

type JobRequestMultiple struct {
	id   string
	exec []*Exec
}

func NewJobRequestMultiple(id string, exec ...*Exec) *JobRequestMultiple {
	return &JobRequestMultiple{
		id:   id,
		exec: exec,
	}
}

func (jr *JobRequestMultiple) AppendExec(exec *Exec) {
	jr.exec = append(jr.exec, exec)
}

func (jr *JobRequestMultiple) SetID(id string) {
	jr.id = id
}

func (jr JobRequestMultiple) GetID() string {
	return jr.id
}

func (jr JobRequestMultiple) GetExec() []*Exec {
	return jr.exec
}
