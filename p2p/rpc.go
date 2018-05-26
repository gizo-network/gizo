package p2p

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/gizo-network/gizo/job/batch"

	"github.com/gizo-network/gizo/job/chain"

	"github.com/gizo-network/gizo/job/chord"

	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job/solo"

	"github.com/gizo-network/gizo/job"
)

func (d Dispatcher) Rpc() {
	d.GetRPC().AddFunction("Version", d.Version)
	d.GetRPC().AddFunction("PeerCount", d.PeerCount)
	d.GetRPC().AddFunction("BlockByHash", d.BlockByHash)
	d.GetRPC().AddFunction("BlockByHeight", d.BlockByHeight)
	d.GetRPC().AddFunction("Latest15Blocks", d.Latest15Blocks)
	d.GetRPC().AddFunction("LatestBlock", d.LatestBlock)
	d.GetRPC().AddFunction("PendingCount", d.PendingCount)
	d.GetRPC().AddFunction("Score", d.Score)
	d.GetRPC().AddFunction("Peers", d.Peers)
	d.GetRPC().AddFunction("PublicKey", d.PublicKey)
	d.GetRPC().AddFunction("NewJob", d.NewJob)
	d.GetRPC().AddFunction("NewExec", d.NewExec)
	d.GetRPC().AddFunction("WorkersCount", d.WorkersCount)
	d.GetRPC().AddFunction("WorkersCountBusy", d.WorkersCountBusy)
	d.GetRPC().AddFunction("WorkersCountNotBusy", d.WorkersCountNotBusy)
	d.GetRPC().AddFunction("ExecStatus", d.ExecStatus)
	d.GetRPC().AddFunction("CancelExec", d.CancelExec)
	d.GetRPC().AddFunction("ExecTimestamp", d.ExecTimestamp)
	d.GetRPC().AddFunction("ExecTimestampString", d.ExecTimestampString)
	d.GetRPC().AddFunction("ExecDurationNanoseconds", d.ExecDurationNanoseconds)
	d.GetRPC().AddFunction("ExecDurationSeconds", d.ExecDurationSeconds)
	d.GetRPC().AddFunction("ExecDurationMinutes", d.ExecDurationMinutes)
	d.GetRPC().AddFunction("ExecDurationString", d.ExecDurationString)
	d.GetRPC().AddFunction("ExecArgs", d.ExecArgs)
	d.GetRPC().AddFunction("ExecErr", d.ExecErr)
	d.GetRPC().AddFunction("ExecPriority", d.ExecPriority)
	d.GetRPC().AddFunction("ExecResult", d.ExecResult)
	d.GetRPC().AddFunction("ExecRetries", d.ExecRetries)
	d.GetRPC().AddFunction("ExecBackoff", d.ExecBackoff)
	d.GetRPC().AddFunction("ExecExecutionTime", d.ExecExecutionTime)
	d.GetRPC().AddFunction("ExecExecutionTimeString", d.ExecExecutionTimeString)
	d.GetRPC().AddFunction("ExecInterval", d.ExecInterval)
	d.GetRPC().AddFunction("ExecBy", d.ExecBy)
	d.GetRPC().AddFunction("ExecTtlNanoseconds", d.ExecTtlNanoseconds)
	d.GetRPC().AddFunction("ExecTtlSeconds", d.ExecTtlSeconds)
	d.GetRPC().AddFunction("ExecTtlMinutes", d.ExecTtlMinutes)
	d.GetRPC().AddFunction("ExecTtlHours", d.ExecTtlHours)
	d.GetRPC().AddFunction("ExecTtlString", d.ExecTtlString)
	d.GetRPC().AddFunction("JobQueueCount", d.JobQueueCount)
	d.GetRPC().AddFunction("LatestBlockHeight", d.LatestBlockHeight)
	d.GetRPC().AddFunction("Job", d.Job)
	d.GetRPC().AddFunction("JobSubmissionTimeUnix", d.JobSubmisstionTimeUnix)
	d.GetRPC().AddFunction("JobSubmissionTimeString", d.JobSubmisstionTimeString)
	d.GetRPC().AddFunction("IsJobPrivate", d.IsJobPrivate)
	d.GetRPC().AddFunction("JobName", d.JobName)
	d.GetRPC().AddFunction("JobLatestExec", d.JobLatestExec)
	d.GetRPC().AddFunction("JobExecs", d.JobExecs)
	d.GetRPC().AddFunction("BlockHashesHex", d.BlockHashesHex)
	d.GetRPC().AddFunction("KeyPair", d.KeyPair)
	d.GetRPC().AddFunction("Solo", d.Solo)
	d.GetRPC().AddFunction("Chord", d.Chord)
	d.GetRPC().AddFunction("Chain", d.Chain)
	d.GetRPC().AddFunction("Batch", d.Batch)
}

//Version returns nodes version
func (d Dispatcher) Version() string {
	return string(NewVersion(GizoVersion, int(d.GetBC().GetLatestHeight()), d.GetBC().GetBlockHashesHex()).Serialize())
}

//PeerCount returns the number of peers a node has
func (d Dispatcher) PeerCount() int {
	return len(d.GetPeers())
}

func (d Dispatcher) BlockByHash(hash string) (string, error) {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return "", err
	}
	b, err := d.GetBC().GetBlockInfo(hashBytes)
	if err != nil {
		return "", err
	}
	return string(b.GetBlock().Serialize()), nil
}

func (d Dispatcher) BlockByHeight(height int) (string, error) {
	b, err := d.GetBC().GetBlockByHeight(height)
	if err != nil {
		return "", err
	}
	return string(b.Serialize()), nil
}

func (d Dispatcher) Latest15Blocks() (string, error) {
	blocksBytes, err := json.Marshal(d.GetBC().GetLatest15())
	if err != nil {
		return "", err
	}
	return string(blocksBytes), nil
}

func (d Dispatcher) LatestBlock() string {
	return string(d.GetBC().GetLatestBlock().Serialize())
}

//PendingCount returns number of job waiting to be written to the bc
func (d Dispatcher) PendingCount() int {
	return d.GetWriteQ().Size()
}

func (d Dispatcher) Score() float64 {
	return d.GetBench().GetScore()
}

//Peers returns the public keys of its peers
func (d Dispatcher) Peers() []string {
	return d.GetPeersPubs()
}

func (d Dispatcher) PublicKey() string {
	return d.GetPubString()
}

func (d Dispatcher) NewJob(task string, name string, priv bool, privKey string) (string, error) {
	if privKey == "" {
		return "", errors.New("Empty private key")
	}
	j, err := job.NewJob(task, name, priv, privKey)
	if err != nil {
		return "", err
	}
	d.AddJob(*j)
	return j.GetID(), nil
}

func (d Dispatcher) NewExec(args []interface{}, retries, priority int, backoff int64, execTime int64, interval int, ttl int64, pub string, envs string) (string, error) {
	e, err := job.DeserializeEnvs([]byte(envs))
	if err != nil {
		return "", err
	}
	exec, err := job.NewExec(args, retries, priority, time.Duration(backoff), execTime, interval, time.Duration(ttl), pub, e, d.GetPubString())
	if err != nil {
		return "", err
	}
	return string(exec.Serialize()), nil
}

func (d Dispatcher) WorkersCount() int {
	return len(d.GetWorkers())
}

func (d Dispatcher) WorkersCountBusy() int {
	temp := 0
	d.mu.Lock()
	for _, info := range d.GetWorkers() {
		if info.GetJob() != nil {
			temp++
		}
	}
	d.mu.Unlock()
	return temp
}

func (d Dispatcher) WorkersCountNotBusy() int {
	temp := 0
	d.mu.Lock()
	for _, info := range d.GetWorkers() {
		if info.GetJob() == nil {
			temp++
		}
	}
	d.mu.Unlock()
	return temp
}

func (d Dispatcher) ExecStatus(id string, hash []byte) (string, error) {
	if d.GetJobPQ().GetPQ().InQueueHash(hash) {
		return job.QUEUED, nil
	} else if worker := d.GetAssignedWorker(hash); worker != nil {
		return job.RUNNING, nil
	}
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetStatus(), nil
}

func (d Dispatcher) CancelExec(hash []byte) error {
	if worker := d.GetAssignedWorker(hash); worker != nil {
		worker.Write(CancelMessage(d.GetPrivByte()))
		return nil
	}
	return errors.New("Exec not running")
}

func (d Dispatcher) ExecTimestamp(id string, hash []byte) (int64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTimestamp(), nil
}

func (d Dispatcher) ExecTimestampString(id string, hash []byte) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return time.Unix(e.GetTimestamp(), 0).String(), nil
}

func (d Dispatcher) ExecDurationNanoseconds(id string, hash []byte) (int64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetDuration().Nanoseconds(), nil
}

func (d Dispatcher) ExecDurationSeconds(id string, hash []byte) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetDuration().Seconds(), nil
}

func (d Dispatcher) ExecDurationMinutes(id string, hash []byte) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetDuration().Minutes(), nil
}

func (d Dispatcher) ExecDurationString(id string, hash []byte) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetDuration().String(), nil
}

func (d Dispatcher) ExecArgs(id string, hash []byte) ([]interface{}, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return make([]interface{}, 1), err
	}
	return e.GetArgs(), nil
}

func (d Dispatcher) ExecErr(id string, hash []byte) (interface{}, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetErr(), nil
}

func (d Dispatcher) ExecPriority(id string, hash []byte) (int, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetPriority(), nil
}

func (d Dispatcher) ExecResult(id string, hash []byte) (interface{}, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetResult(), nil
}

func (d Dispatcher) ExecRetries(id string, hash []byte) (int, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetRetries(), nil
}

func (d Dispatcher) ExecBackoff(id string, hash []byte) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetBackoff().Seconds(), nil
}

func (d Dispatcher) ExecExecutionTime(id string, hash []byte) (int64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetExecutionTime(), nil
}

func (d Dispatcher) ExecExecutionTimeString(id string, hash []byte) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return time.Unix(e.GetExecutionTime(), 0).String(), nil
}

func (d Dispatcher) ExecInterval(id string, hash []byte) (int, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetInterval(), nil
}

func (d Dispatcher) ExecBy(id string, hash []byte) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetBy(), nil
}

func (d Dispatcher) ExecTtlNanoseconds(id string, hash []byte) (int64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTTL().Nanoseconds(), nil
}

func (d Dispatcher) ExecTtlSeconds(id string, hash []byte) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTTL().Seconds(), nil
}

func (d Dispatcher) ExecTtlMinutes(id string, hash []byte) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTTL().Minutes(), nil
}

func (d Dispatcher) ExecTtlHours(id string, hash []byte) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTTL().Hours(), nil
}

func (d Dispatcher) ExecTtlString(id string, hash []byte) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetTTL().String(), nil
}

//JobQueueCount returns nubmer of jobs waiting to be executed
func (d Dispatcher) JobQueueCount() int {
	return d.GetJobPQ().Len()
}
func (d Dispatcher) LatestBlockHeight() int {
	return int(d.GetBC().GetLatestBlock().GetHeight())
}

//Job returns a job
func (d Dispatcher) Job(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return string(j.Serialize()), nil
}

func (d Dispatcher) JobSubmisstionTimeUnix(id string) (int64, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return 0, err
	}
	return j.GetSubmissionTime().Unix(), nil
}

func (d Dispatcher) JobSubmisstionTimeString(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return j.GetSubmissionTime().String(), nil
}

func (d Dispatcher) IsJobPrivate(id string) (bool, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return false, err
	}
	return j.GetPrivate(), nil
}

func (d Dispatcher) JobName(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return j.GetName(), nil
}

func (d Dispatcher) JobLatestExec(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return string(j.GetLatestExec().Serialize()), nil
}

func (d Dispatcher) JobExecs(id string) (string, error) {
	execsBytes, err := json.Marshal(d.GetBC().GetJobExecs(id))
	if err != nil {
		return "", err
	}
	return string(execsBytes), nil
}

//BlockHashesHex returns hashes of all blocks in the bc
func (d Dispatcher) BlockHashesHex() []string {
	return d.GetBC().GetBlockHashesHex()
}

//KeyPair returns new pub and priv keypair
func (d Dispatcher) KeyPair() (string, error) {
	priv, pub := crypt.GenKeys()
	temp := make(map[string]string)
	temp["priv"] = hex.EncodeToString(priv)
	temp["pub"] = hex.EncodeToString(pub)
	keysBytes, err := json.Marshal(temp)
	if err != nil {
		return "", err
	}
	return string(keysBytes), nil
}

func (d Dispatcher) Solo(jr string) (string, error) {
	//TODO: send result to message broker
	request, err := job.DeserializeJRS([]byte(jr))
	if err != nil {
		return "", err
	}
	solo := solo.NewSolo(request, d.GetBC(), d.GetJobPQ(), d.GetJC())
	solo.Dispatch() //FIXME: look for a more controllable solution
	return string(solo.Result().Serialize()), nil
}

func (d Dispatcher) Chord(jrs []string, callbackJr string) (string, error) {
	//TODO: send result to message broker
	var requests []job.JobRequestMultiple
	for _, jr := range jrs {
		request, err := job.DeserializeJRM([]byte(jr))
		if err != nil {
			return "", err
		}
		requests = append(requests, request)
	}
	callackRequest, err := job.DeserializeJRM([]byte(callbackJr))
	if err != nil {
		return "", err
	}
	c, err := chord.NewChord(requests, callackRequest, d.GetBC(), d.GetJobPQ(), d.GetJC())
	if err != nil {
		return "", err
	}
	c.Dispatch()
	return string(c.Result().Serialize()), nil
}

func (d Dispatcher) Chain(jrs []string, callbackJr string) (string, error) {
	//TODO: send result to message broker
	var requests []job.JobRequestMultiple
	for _, jr := range jrs {
		request, err := job.DeserializeJRM([]byte(jr))
		if err != nil {
			return "", err
		}
		requests = append(requests, request)
	}

	c, err := chain.NewChain(requests, d.GetBC(), d.GetJobPQ(), d.GetJC())
	if err != nil {
		return "", err
	}
	c.Dispatch()
	result, err := json.Marshal(c.Result())
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (d Dispatcher) Batch(jrs []string, callbackJr string) (string, error) {
	//TODO: send result to message broker
	var requests []job.JobRequestMultiple
	for _, jr := range jrs {
		request, err := job.DeserializeJRM([]byte(jr))
		if err != nil {
			return "", err
		}
		requests = append(requests, request)
	}

	b, err := batch.NewBatch(requests, d.GetBC(), d.GetJobPQ(), d.GetJC())
	if err != nil {
		return "", err
	}
	b.Dispatch()
	result, err := json.Marshal(b.Result())
	if err != nil {
		return "", err
	}
	return string(result), nil
}
