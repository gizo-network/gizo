package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
)

func TestDifficulty(t *testing.T) {
	core.RemoveDataPath()
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})

	nodes := []*merkletree.MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}
	tree := merkletree.NewMerkleTree(nodes)
	block := core.NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 3, 10)
	bc := core.CreateBlockChain()
	bc.AddBlock(block)
	b := benchmark.NewBenchmarkEngine()
	assert.NotNil(t, Difficulty(b.GetData(), *bc))
}
