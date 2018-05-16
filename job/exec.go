package job

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/kpango/glg"
)

//TODO: add environment variables
type Exec struct {
	Hash          []byte
	Timestamp     int64
	Duration      time.Duration //saved in nanoseconds
	Args          []interface{} // parameters
	Err           interface{}
	Priority      int
	Result        interface{}
	Status        string        //job status
	Retries       int           // number of max retries
	RetriesCount  int           //number of retries
	Backoff       time.Duration //backoff time of retries (seconds)
	ExecutionTime int64         // time scheduled to run (unix) - should sleep # of seconds before adding to job queue
	Interval      int           //periodic job exec (seconds)
	By            string        //! ID of the worker node that ran this
	TTL           time.Duration //! time limit of job running
	Pub           string        //! public key for private jobs
	Envs          []byte
	cancel        chan struct{}
}

func NewExec(args []interface{}, retries, priority int, backoff time.Duration, execTime int64, interval int, ttl time.Duration, pub string, envs EnvironmentVariables, passphrase string) (*Exec, error) {
	if retries > MaxRetries {
		return nil, ErrRetriesOutsideLimit
	}

	encryptEnvs := helpers.Encrypt(envs.Serialize(), passphrase)
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
		Envs:          encryptEnvs,
		Pub:           pub,
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

func (e Exec) GetEnvs(passphrase string) (EnvironmentVariables, error) {
	d, err := helpers.Decrypt(e.Envs, passphrase)
	if err != nil {
		return EnvironmentVariables{}, errors.New("Unable to decrypt environment variables")
	}
	return DeserializeEnvs(d)
}

//returns environment variables as a map
func (e Exec) GetEnvsMap(passphrase string) (map[string]interface{}, error) {
	temp := make(map[string]interface{})
	envs, err := e.GetEnvs(passphrase)
	if err != nil {
		return nil, err
	}
	for _, val := range envs {
		temp[val.GetKey()] = val.GetValue()
	}
	return temp, nil
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
			[]byte(e.GetBy()),
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

func (e Exec) GetBy() string {
	return e.By
}

func (e *Exec) SetBy(by string) {
	e.By = by
}

func (e Exec) getPub() string {
	return e.Pub
}

func (e Exec) Serialize() []byte {
	temp, err := json.Marshal(e)
	if err != nil {
		glg.Error(err)
	}
	return temp
}

func DeserializeExec(b []byte) Exec {
	var temp Exec
	err := json.Unmarshal(b, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	temp.cancel = make(chan struct{})
	return temp
}

//UniqExec returns unique values of parameter
func UniqExec(execs []Exec) []Exec {
	temp := []Exec{}
	seen := make(map[string]bool)
	for _, exec := range execs {
		if _, ok := seen[string(exec.Serialize())]; ok {
			continue
		}
		seen[string(exec.Serialize())] = true
		temp = append(temp, exec)
	}
	return temp
}
