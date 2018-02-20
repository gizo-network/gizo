package job

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"strconv"
	"time"

	"github.com/kpango/glg"
)

type JobExec struct {
	Hash          []byte        `json:"hash"`
	Timestamp     int64         `json:"timestamp"`
	Duration      time.Duration `json:"duaration"` //saved in nanoseconds
	Args          []interface{} `json:"args"`
	Err           interface{}   `json:"err"`
	Result        interface{}   `json:"result"`
	Status        string        `json:"status"`         //job status
	Retries       int           `json:"retries"`        // number of retries
	RetryDelay    time.Duration `json:"retry_delay"`    //backoff time of retries (seconds)
	ExecutionTime time.Duration `json:"execution_time"` // time scheduled to run (seconds)
	By            []byte        `json:"by"`             //! ID of the worker node that ran this
}

func (j JobExec) GetArgs() []interface{} {
	return j.Args
}

func (j *JobExec) SetArgs(a []interface{}) {
	j.Args = a
}

func (j JobExec) GetHash() []byte {
	return j.Hash
}

func (j *JobExec) setHash() {
	e, err := json.Marshal(j.GetErr())
	if err != nil {
		glg.Error(err)
	}
	result, err := json.Marshal(j.GetResult())
	if err != nil {
		glg.Error(err)
	}

	header := bytes.Join(
		[][]byte{
			[]byte(strconv.FormatInt(j.GetTimestamp(), 10)),
			[]byte(strconv.FormatInt(int64(j.GetDuration()), 10)),
			e,
			result,
			j.GetBy(),
		},
		[]byte{},
	)

	hash := sha256.Sum256(header)
	j.Hash = hash[:]
}

func (j JobExec) GetTimestamp() int64 {
	return j.Timestamp
}

func (j *JobExec) SetTimestamp(t int64) {
	j.Timestamp = t
}

func (j JobExec) GetDuration() time.Duration {
	return j.Duration
}

func (j *JobExec) SetDuration(t time.Duration) {
	j.Duration = t
}

func (j JobExec) GetErr() interface{} {
	return j.Err
}

func (j *JobExec) SetErr(e interface{}) {
	j.Err = e
}

func (j JobExec) GetResult() interface{} {
	return j.Result
}

func (j *JobExec) SetResult(r interface{}) {
	j.Result = r
}

func (j JobExec) GetBy() []byte {
	return j.By
}

func (j *JobExec) SetBy(by []byte) {
	j.By = by
}

func (j JobExec) Serialize() []byte {
	temp, err := json.Marshal(j)
	if err != nil {
		glg.Error(err)
	}
	return temp
}
