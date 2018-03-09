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
}

func NewExec(args []interface{}, retries, priority int, backoff time.Duration, execTime int64, interval int, ttl time.Duration, pub string, envs EnvironmentVariables) (*Exec, error) {
	if retries > MaxRetries {
		return nil, ErrRetriesOutsideLimit
	}
	return &Exec{
		Args:          args,
		Retries:       retries,
		RetriesCount:  0,
		Priority:      priority,
		Status:        STARTED,
		Backoff:       backoff,
		ExecutionTime: execTime,
		Interval:      interval,
		TTL:           ttl,
		envs:          envs,
		By:            []byte("0000"), //!FIXME: replace with real node ID
		pub:           pub,
	}, nil
}

func (j Exec) GetEnvs() EnvironmentVariables {
	return j.envs
}

func (j Exec) GetEnvsMap() map[string]interface{} {
	temp := make(map[string]interface{})
	for _, val := range j.GetEnvs() {
		temp[val.GetKey()] = val.GetValue()
	}
	return temp
}

func (j Exec) GetTTL() time.Duration {
	return j.TTL
}

func (j *Exec) SetTTL(ttl time.Duration) {
	j.TTL = ttl
}

func (j Exec) GetInterval() int {
	return j.Interval
}

func (j *Exec) SetInterval(i int) {
	j.Interval = i
}

func (j Exec) GetPriority() int {
	return j.Priority
}

func (j *Exec) SetPriority(p int) error {
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

func (j Exec) GetExecutionTime() int64 {
	return j.ExecutionTime
}

//? takes unix time
func (j *Exec) SetExecutionTime(e int64) error {
	if time.Unix(e, 0).Before(time.Now()) {
		return ErrExecutionTimeBehind
	}
	j.ExecutionTime = e
	return nil
}

func (j Exec) GetBackoff() time.Duration {
	return j.Backoff
}

func (j *Exec) SetBackoff(b time.Duration) error {
	if b > MaxRetryBackoff {
		return ErrRetryDelayOutsideLimit
	}
	j.Backoff = b
	return nil
}

func (j Exec) GetRetriesCount() int {
	return j.RetriesCount
}

func (j *Exec) IncrRetriesCount() {
	j.RetriesCount++
}

func (j Exec) GetRetries() int {
	return j.Retries
}

func (j *Exec) SetRetries(r int) error {
	if r > MaxRetries {
		return ErrRetriesOutsideLimit
	}
	j.Retries = r
	return nil
}

func (j Exec) GetStatus() string {
	return j.Status
}

func (j *Exec) SetStatus(s string) {
	j.Status = s
}

func (j Exec) GetArgs() []interface{} {
	return j.Args
}

func (j *Exec) SetArgs(a []interface{}) {
	j.Args = a
}

func (j Exec) GetHash() []byte {
	return j.Hash
}

func (j *Exec) setHash() {
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

func (j Exec) GetTimestamp() int64 {
	return j.Timestamp
}

func (j *Exec) SetTimestamp(t int64) {
	j.Timestamp = t
}

func (j Exec) GetDuration() time.Duration {
	return j.Duration
}

func (j *Exec) SetDuration(t time.Duration) {
	j.Duration = t
}

func (j Exec) GetErr() interface{} {
	return j.Err
}

func (j *Exec) SetErr(e interface{}) {
	j.Err = e
}

func (j Exec) GetResult() interface{} {
	return j.Result
}

func (j *Exec) SetResult(r interface{}) {
	j.Result = r
}

func (j Exec) GetBy() []byte {
	return j.By
}

func (j *Exec) SetBy(by []byte) {
	j.By = by
}

func (j Exec) Serialize() []byte {
	temp, err := json.Marshal(j)
	if err != nil {
		glg.Error(err)
	}
	return temp
}

func (j Exec) getPub() string {
	return j.pub
}
