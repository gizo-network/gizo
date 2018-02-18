package job

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kpango/glg"
	"github.com/mattn/anko/vm"
)

var (
	ErrExecNotFound = errors.New("Exec Not Found")
)

type Job struct {
	ID        uuid.UUID `json:"id"`
	Hash      []byte    `json:"hash"`
	Exec      []JobExec `json:"jobexec"`
	Source    []byte    `json:"source"`
	Signature []byte    `json:"signature"` // signature of deployer
}

func NewJob(s string) *Job {
	j := &Job{
		ID:     uuid.New(),
		Exec:   []JobExec{},
		Source: []byte(s),
	}
	j.setHash()
	return j
}

func (j Job) GetID() uuid.UUID {
	return j.ID
}

func (j Job) GetHash() []byte {
	return j.Hash
}

func (j *Job) setHash() {
	headers := bytes.Join(
		[][]byte{
			[]byte(j.GetID().String()),
			j.serializeExecs(),
			j.GetSource(),
		},
		[]byte{},
	)
	hash := sha256.Sum256(headers)
	j.Hash = hash[:]
}

func (j Job) serializeExecs() []byte {
	temp, err := json.Marshal(j.GetExecs())
	if err != nil {
		glg.Error(err)
	}
	return temp
}

func (j Job) GetExec(hash []byte) (*JobExec, error) {
	var check int
	for _, exec := range j.GetExecs() {
		check = bytes.Compare(exec.GetHash(), hash)
		if check == 0 {
			return &exec, nil
		}
	}
	return nil, ErrExecNotFound
}

func (j Job) GetLatestExec() JobExec {
	return j.Exec[len(j.GetExecs())-1]
}

func (j Job) GetExecs() []JobExec {
	return j.Exec
}

func (j *Job) AddExec(je JobExec) {
	j.Exec = append(j.Exec, je)
	j.setHash() //regenerates hash
}

func (j Job) GetSource() []byte {
	return j.Source
}

func (j *Job) SetSource(s []byte) {
	j.Source = s
}

func (j *Job) Serialize() []byte {
	temp, err := json.Marshal(*j)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}

//FIXME: add fault tolerance and security
func (j *Job) Execute() (interface{}, error) {
	env := vm.NewEnv()
	start := time.Now()
	result, err := env.Execute(string(j.GetSource()))
	exec := JobExec{
		Timestamp: time.Now().Unix(),
		Duration:  time.Now().Sub(start).Nanoseconds(),
		Err:       err,
		Result:    result,
		By:        []byte("0000"), //FIXME: replace with real ID
	}
	j.AddExec(exec)
	return result, err
}
