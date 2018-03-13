package job

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"strconv"
	"time"

	"github.com/kpango/glg"
)

//TODO: add environment variables
type Exec struct {
	Hash          []byte        `json:"hash"`
	Timestamp     int64         `json:"timestamp"`
	Duration      time.Duration `json:"duaration"` //saved in nanoseconds
	Args          []interface{} `json:"args"`
	Err           interface{}   `json:"err"`
	Priority      int           `json:"priority"`
	Result        interface{}   `json:"result"`
	Status        string        `json:"status"`         //job status
	Retries       int           `json:"retries"`        // number of max retries
	RetriesCount  int           `json:"retries_count"`  //number of retries
	Backoff       time.Duration `json:"backoff"`        //backoff time of retries (seconds)
	ExecutionTime int64         `json:"execution_time"` // time scheduled to run (unix) - should sleep # of seconds before adding to job queue
	Interval      int           `json:"interval"`       //periodic job exec (seconds)
	By            []byte        `json:"by"`             //! ID of the worker node that ran this
	TTL           time.Duration `json:"ttl"`            //! time limit of job running
	envs          EnvironmentVariables
	pub           string //! public key for private jobs
	cancel        chan struct{}
}

func NewExec(args []interface{}, retries, priority int, backoff time.Duration, execTime int64, interval int, ttl time.Duration, pub string, envs EnvironmentVariables) (*Exec, error) {
	if retries > MaxRetries {
		return nil, ErrRetriesOutsideLimit
	}

	ex := &Exec{
		Args:          args,
		Retries:       retries,
		RetriesCount:  0, //initialized to 0
		Priority:      priority,
		Status:        STARTED,
		Backoff:       backoff,
		Interval:      interval,
		ExecutionTime: execTime,
		TTL:           ttl,
		envs:          envs,
		By:            []byte("0000"), //!FIXME: replace with real node ID
		pub:           pub,
		cancel:        make(chan struct{}),
	}
	return ex, nil
}

func (e *Exec) Cancel() {
	e.cancel <- struct{}{}
}

func (e Exec) GetCancelChan() chan struct{} {
	return e.cancel
}

func (e Exec) GetEnvs() EnvironmentVariables {
	return e.envs
}

func (e Exec) GetEnvsMap() map[string]interface{} {
	temp := make(map[string]interface{})
	for _, val := range e.GetEnvs() {
		temp[val.GetKey()] = val.GetValue()
	}
	return temp
}

func (e Exec) GetTTL() time.Duration {
	return e.TTL
}

func (e *Exec) SetTTL(ttl time.Duration) {
	e.TTL = ttl
}

func (e Exec) GetInterval() int {
	return e.Interval
}

func (e *Exec) SetInterval(i int) {
	e.Interval = i
}

func (e Exec) GetPriority() int {
	return e.Priority
}

func (e *Exec) SetPriority(p int) error {
	switch p {
	case HIGH:
	case MEDIUM:
	case LOW:
	case NORMAL:
	default:
		return ErrInvalidPriority
	}
	return nil
}

func (e Exec) GetExecutionTime() int64 {
	return e.ExecutionTime
}

//? takes unix time
func (e *Exec) SetExecutionTime(t int64) error {
	if time.Now().Unix() > t {
		return ErrExecutionTimeBehind
	}
	e.ExecutionTime = t
	return nil
}

func (e Exec) GetBackoff() time.Duration {
	return e.Backoff
}

func (e *Exec) SetBackoff(b time.Duration) error {
	if b > MaxRetryBackoff {
		return ErrRetryDelayOutsideLimit
	}
	e.Backoff = b
	return nil
}

func (e Exec) GetRetriesCount() int {
	return e.RetriesCount
}

func (e *Exec) IncrRetriesCount() {
	e.RetriesCount++
}

func (e Exec) GetRetries() int {
	return e.Retries
}

func (e *Exec) SetRetries(r int) error {
	if r > MaxRetries {
		return ErrRetriesOutsideLimit
	}
	e.Retries = r
	return nil
}

func (e Exec) GetStatus() string {
	return e.Status
}

func (e *Exec) SetStatus(s string) {
	e.Status = s
}

func (e Exec) GetArgs() []interface{} {
	return e.Args
}

func (e *Exec) SetArgs(a []interface{}) {
	e.Args = a
}

func (e Exec) GetHash() []byte {
	return e.Hash
}

func (e *Exec) setHash() {
	stringified, err := json.Marshal(e.GetErr())
	if err != nil {
		glg.Error(err)
	}
	result, err := json.Marshal(e.GetResult())
	if err != nil {
		glg.Error(err)
	}

	header := bytes.Join(
		[][]byte{
			[]byte(strconv.FormatInt(e.GetTimestamp(), 10)),
			[]byte(strconv.FormatInt(int64(e.GetDuration()), 10)),
			stringified,
			result,
			e.GetBy(),
		},
		[]byte{},
	)

	hash := sha256.Sum256(header)
	e.Hash = hash[:]
}

func (e Exec) GetTimestamp() int64 {
	return e.Timestamp
}

func (e *Exec) SetTimestamp(t int64) {
	e.Timestamp = t
}

func (e Exec) GetDuration() time.Duration {
	return e.Duration
}

func (e *Exec) SetDuration(t time.Duration) {
	e.Duration = t
}

func (e Exec) GetErr() interface{} {
	return e.Err
}

func (e *Exec) SetErr(err interface{}) {
	e.Err = err
}

func (e Exec) GetResult() interface{} {
	return e.Result
}

func (e *Exec) SetResult(r interface{}) {
	e.Result = r
}

func (e Exec) GetBy() []byte {
	return e.By
}

func (e *Exec) SetBy(by []byte) {
	e.By = by
}

func (e Exec) Serialize() []byte {
	temp, err := json.Marshal(e)
	if err != nil {
		glg.Error(err)
	}
	return temp
}

func (e Exec) getPub() string {
	return e.pub
}
