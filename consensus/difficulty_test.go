package consensus_test

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/consensus"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
)

func TestDifficulty(t *testing.T) {
	godotenv.Load()
	core.RemoveDataPath()
	bc := core.CreateBlockChain()
	j := job.NewJob("func test(){return 1+1}", "test")
	node1 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})

	nodes := []*merkletree.MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}
	tree := merkletree.NewMerkleTree(nodes)
	block := core.NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestHeight(), 10)
	bc.AddBlock(block)
	b := benchmark.NewBenchmarkEngine()
	assert.NotNil(t, consensus.Difficulty(b.GetData(), *bc))
}
