package batch_test

import (
	"testing"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/batch"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	core.RemoveDataPath()
	pq := queue.NewJobPriorityQueue()
	j := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial")
	j2 := job.NewJob(`
		func Test(){
			return "Testing"
		}
		`, "Test")
	exec1, err := job.NewExec([]interface{}{"10"}, 5, job.NORMAL, 0, 0, 0, 0)
	assert.NoError(t, err)
	exec2, err := job.NewExec([]interface{}{"11"}, 5, job.NORMAL, 0, 0, 0, 0)
	assert.NoError(t, err)
	exec3, err := job.NewExec([]interface{}{"12"}, 5, job.NORMAL, 0, 0, 0, 0)
	assert.NoError(t, err)
	exec4, err := job.NewExec([]interface{}{""}, 5, job.NORMAL, 0, 0, 0, 0)
	assert.NoError(t, err)
	node1 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode(*j2, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	nodes := []*merkletree.MerkleNode{node1, node2}
	tree := merkletree.NewMerkleTree(nodes)
	bc := core.CreateBlockChain()
	block := core.NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestHeight()+1, 10)
	bc.AddBlock(block)
	jr := job.NewJobRequest(j.GetID(), exec1, exec2, exec3)
	jr2 := job.NewJobRequest(j2.GetID(), exec4, exec4, exec4, exec4, exec4)
	batch, err := batch.NewBatch([]job.JobRequest{*jr, *jr2}, bc, pq)
	assert.NoError(t, err)
	batch.Dispatch()
	assert.NotNil(t, batch.Result())
}
