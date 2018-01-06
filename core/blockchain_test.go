package core

import (
	"testing"

	"github.com/gizo-network/gizo/core/merkle_tree"
	"github.com/stretchr/testify/assert"
)

func TestNewBlockChain(t *testing.T) {
	bc := NewBlockChain()
	assert.NotNil(t, bc)
	assert.NotNil(t, bc.Blocks)
}

func TestAddBlock(t *testing.T) {
	node1 := merkle_tree.NewNode([]byte("test1asdfasdf job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node2 := merkle_tree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node3 := merkle_tree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node4 := merkle_tree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	nodes := []*merkle_tree.MerkleNode{node1, node2, node3, node4}
	tree := merkle_tree.NewMerkleTree(nodes)
	bc := NewBlockChain()
	bc.AddBlock(*tree)
	assert.NotEmpty(t, bc.Blocks, "empty block")
	assert.Equal(t, 2, len(bc.Blocks), "chain height not 2")
}

func TestVerifyBlockChain(t *testing.T) {
	node1 := merkle_tree.NewNode([]byte("test1asdfasdf job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node2 := merkle_tree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node3 := merkle_tree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node4 := merkle_tree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node5 := merkle_tree.NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node6 := merkle_tree.NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node7 := merkle_tree.NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node8 := merkle_tree.NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})

	tree1 := merkle_tree.NewMerkleTree([]*merkle_tree.MerkleNode{node1, node2, node3, node4})
	tree2 := merkle_tree.NewMerkleTree([]*merkle_tree.MerkleNode{node5, node6, node7, node8})

	blockchain := NewBlockChain()
	blockchain.AddBlock(*tree1)
	blockchain.AddBlock(*tree2)
	assert.True(t, blockchain.VerifyBlockChain(), "blockchain not verified")

	//modify a single value
	blockchain.Blocks[1].Nonce = 40
	assert.False(t, blockchain.VerifyBlockChain(), "blockchain verified")
}
