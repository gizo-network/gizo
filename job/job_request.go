package job

type JobRequest struct {
	id   string
	exec []*Exec
}

func NewJobRequest(id string, exec ...*Exec) *JobRequest {
	return &JobRequest{
		id:   id,
		exec: exec,
	}
}

func (jr JobRequest) GetID() string {
	return jr.id
}

func (jr JobRequest) GetExec() []*Exec {
	return jr.exec
}
