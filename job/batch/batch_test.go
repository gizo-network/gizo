package batch_test

import (
	"encoding/hex"
	"testing"

	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/batch"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	core.RemoveDataPath()
	priv, _ := crypt.GenKeys()
	pq := queue.NewJobPriorityQueue()
	j := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, hex.EncodeToString(priv))
	j2 := job.NewJob(`
		func Test(){
			return "Testing"
		}
		`, "Test", false, hex.EncodeToString(priv))
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec1, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs)
	assert.NoError(t, err)
	exec2, err := job.NewExec([]interface{}{11}, 5, job.NORMAL, 0, 0, 0, 0, "", envs)
	assert.NoError(t, err)
	exec3, err := job.NewExec([]interface{}{12}, 5, job.NORMAL, 0, 0, 0, 0, "", envs)
	assert.NoError(t, err)
	exec4, err := job.NewExec([]interface{}{}, 5, job.NORMAL, 0, 0, 0, 0, "", envs)
	assert.NoError(t, err)
	node1 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode(*j2, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	nodes := []*merkletree.MerkleNode{node1, node2}
	tree := merkletree.NewMerkleTree(nodes)
	bc := core.CreateBlockChain("test")
	block := core.NewBlock(*tree, bc.GetPrevHash(), bc.GetNextHeight(), 10, "test")
	bc.AddBlock(block)
	jr := job.NewJobRequestMultiple(j.GetID(), exec1, exec2, exec3)
	jr2 := job.NewJobRequestMultiple(j2.GetID(), exec4, exec4, exec4, exec4, exec4)
	batch, err := batch.NewBatch([]job.JobRequestMultiple{*jr, *jr2}, bc, pq, cache.NewJobCacheNoWatch(bc))
	assert.NoError(t, err)
	batch.Dispatch()
	assert.NotNil(t, batch.Result())
}
