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
	Priority      int           `json:"priority"`
	Result        interface{}   `json:"result"`
	Status        string        `json:"status"`         //job status
	Retries       int           `json:"retries"`        // number of retries
	Backoff       time.Duration `json:"backoff"`        //backoff time of retries (seconds)
	ExecutionTime int64         `json:"execution_time"` // time scheduled to run (unix) - should sleep # of seconds before adding to job queue
	Interval      int           `json:"interval"`       //periodic job exec (seconds)
	By            []byte        `json:"by"`             //! ID of the worker node that ran this
}

func NewJobExec(args []interface{}, retries, priority int, backoff time.Duration, execTime int64, interval int) (*JobExec, error) {
	switch priority {
	case HIGH:
	case MEDIUM:
	case LOW:
	case NORMAL:

	default:
		return nil, ErrInvalidPriority
	}
	if retries > MaxRetries {
		return nil, ErrRetriesOutsideLimit
	}
	return &JobExec{
		Args:          args,
		Retries:       retries,
		Priority:      priority,
		Status:        STARTED,
		Backoff:       backoff,
		ExecutionTime: execTime,
		Interval:      interval,
		By:            []byte("0000"), //!FIXME: replace with real ID
	}, nil
}

func (j JobExec) GetInterval() int {
	return j.Interval
}

func (j *JobExec) SetInterval(i int) {
	j.Interval = i
}

func (j JobExec) GetPriority() int {
	return j.Priority
}

func (j *JobExec) SetPriority(p int) error {
	if p != HIGH || p != MEDIUM || p != LOW || p != NORMAL {
		return ErrInvalidPriority
	}
	return nil
}

func (j JobExec) GetExecutionTime() int64 {
	return j.ExecutionTime
}

//? takes unix time
func (j *JobExec) SetExecutionTime(e int64) error {
	if e < time.Now().Unix() {
		return ErrExecutionTimeBehind
	}
	j.ExecutionTime = e
	return nil
}

func (j JobExec) GetBackoff() time.Duration {
	return j.Backoff
}

func (j *JobExec) SetBackoff(b time.Duration) error {
	if b > MaxRetryDelay {
		return ErrRetryDelayOutsideLimit
	}
	j.Backoff = b
	return nil
}

func (j JobExec) GetRetries() int {
	return j.Retries
}

func (j *JobExec) SetRetries(r int) error {
	if r > MaxRetries {
		return ErrRetriesOutsideLimit
	}
	j.Retries = r
	return nil
}

func (j JobExec) GetStatus() string {
	return j.Status
}

func (j *JobExec) SetStatus(s string) {
	j.Status = s
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
