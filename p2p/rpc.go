package p2p

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/gizo-network/gizo/crypt"

	"github.com/gizo-network/gizo/job"
)

func (d Dispatcher) RpcHttp() {
	d.GetRPC().AddFunction("PeerCount", d.PeerCount)
	d.GetRPC().AddFunction("BlockByHash", d.BlockByHash)
	d.GetRPC().AddFunction("BlockByHeight", d.BlockByHeight)
	d.GetRPC().AddFunction("Latest15Blocks", d.Latest15Blocks)
	d.GetRPC().AddFunction("LatestBlock", d.LatestBlock)
}

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

func (d Dispatcher) PendingCount() int {
	return d.GetWriteQ().Size()
}

func (d Dispatcher) Score() float64 {
	return d.GetBench().GetScore()
}

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
	e := job.DeserializeEnvs([]byte(envs))
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

//TODO: implemented worker cancel and getstatus
// func (d Dispatcher) ExecStatus(id string, hash []byte) (string, error) {
// 	if d.GetJobPQ().GetPQ().InQueueHash(hash) {
// 		return "", job.QUEUED
// 	} else if worker := d.GetAssignedWorker(hash); worker != nil {

// 	} else {
// 		e := d.GetBC().find
// 	}
// }

// func (d Dispatcher) CancelExec() error {

// }

// func (d Dispatcher) GetExecTimestamp() {

// }
// func (d Dispatcher) GetExecTimestampString() {

// }
// func (d Dispatcher) GetExecDurationNanoseconds() {

// }
// GetExecDurationMilliseconds
// GetExecDurationSeconds
// GetExecDurationMinutes
// GetExecDurationHours
// GetExecDurationString
// GetExecArgs
// GetExecErr
// GetExecPriority
// GetExecResult
// GetExecStatus
// GetExecRetries
// GetExecRetriesCount
// GetExecBackoff
// GetExecExcutionTimeUnix
// GetExecExcutionTimeString
// func (d Dispatcher) GetExecIntervalNanoseconds(hash []byte) {

// }
// GetExecIntervalMilliseconds
// GetExecIntervalSeconds
// GetExecIntervalMinutes
// GetExecBy
// GetExecTtlNanoseconds
// GetExecTtlMilliseconds
// GetExecTtlSeconds
// GetExecTtlMinutes
// GetExecTtlHours
// GetExecTtlString
// GetExecPub
// GetExecEnvs
func (d Dispatcher) JobQueueCount() int {
	return d.GetJobPQ().Len()
}
func (d Dispatcher) GetLatestBlockHeight() int {
	return int(d.GetBC().GetLatestBlock().GetHeight())
}

func (d Dispatcher) GetJob(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return string(j.Serialize()), nil
}

func (d Dispatcher) GetJobSubmisstionTimeUnix(id string) (int64, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return 0, err
	}
	return j.GetSubmissionTime().Unix(), nil
}

func (d Dispatcher) GetJobSubmisstionTimeString(id string) (string, error) {
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

func (d Dispatcher) GetJobName(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return j.GetName(), nil
}

func (d Dispatcher) GetJobLatestExec(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return string(j.GetLatestExec().Serialize()), nil
}

func (d Dispatcher) GetJobExecs(id string) (string, error) {
	execsBytes, err := json.Marshal(d.GetBC().GetJobExecs(id))
	if err != nil {
		return "", err
	}
	return string(execsBytes), nil
}

func (d Dispatcher) GetBlockHashesHex() []string {
	return d.GetBlockHashesHex()
}

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

// NewSolo
// NewChord
// NewChain
// NewBatch
