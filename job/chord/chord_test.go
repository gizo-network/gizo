package chord_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/chord"
	"github.com/gizo-network/gizo/job/queue"
)

func TestChord(t *testing.T) {
	core.RemoveDataPath()
	priv, pub := crypt.GenKeys()
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
			return "test"
		}
		`, "Test", false, hex.EncodeToString(priv))
	callback := job.NewJob(`
			func Callback(n, n2, n3, n4){
				println(n, n2, n3, n4)
				return n
			}
			`, "Callback", false, hex.EncodeToString(priv))
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec1, err := job.NewExec([]interface{}{2}, 5, job.NORMAL, 0, 0, 0, 0, hex.EncodeToString(pub), envs)
	assert.NoError(t, err)
	exec2, err := job.NewExec([]interface{}{4}, 5, job.HIGH, 0, 0, 0, 0, hex.EncodeToString(pub), envs)
	assert.NoError(t, err)
	exec3, err := job.NewExec([]interface{}{3}, 5, job.MEDIUM, 0, 0, 0, 0, hex.EncodeToString(pub), envs)
	assert.NoError(t, err)
	exec4, err := job.NewExec([]interface{}{}, 5, job.LOW, 0, 0, 0, 0, "", envs)
	assert.NoError(t, err)
	exec5, err := job.NewExec([]interface{}{}, 5, job.LOW, 0, 0, 0, 0, "", envs)
	assert.NoError(t, err)

	node1 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode(*j2, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode(*callback, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})

	nodes := []*merkletree.MerkleNode{node1, node2, node3}
	tree := merkletree.NewMerkleTree(nodes)
	bc := core.CreateBlockChain("test")
	block := core.NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestHeight()+1, 10)
	bc.AddBlock(block)
	jr := job.NewJobRequestMultiple(j.GetID(), exec1, exec2, exec3)
	jr2 := job.NewJobRequestMultiple(j2.GetID(), exec4)
	callbackJR := job.NewJobRequestMultiple(callback.GetID(), exec5)
	c, err := chord.NewChord([]job.JobRequestMultiple{*jr, *jr2}, *callbackJR, bc, pq, cache.NewJobCacheNoWatch(bc))
	assert.NoError(t, err)
	c.Dispatch()
	assert.NotNil(t, c.Result().GetExec()[0].GetResult())
}
