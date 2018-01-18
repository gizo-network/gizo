package core

// import (
// 	"testing"

// 	"github.com/gizo-network/gizo/core/merkletree"
// 	"github.com/stretchr/testify/assert"
// )

// func TestNewBlockChain(t *testing.T) {
// 	bc := NewBlockChain()
// 	assert.NotNil(t, bc)
// 	assert.NotNil(t, bc.Blocks)
// }

// func TestAddBlock(t *testing.T) {
// 	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node2 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node3 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node4 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	nodes := []*merkletree.MerkleNode{node1, node2, node3, node4}
// 	tree := merkletree.NewMerkleTree(nodes)
// 	bc := NewBlockChain()
// 	bc.AddBlock(*tree)
// 	assert.NotEmpty(t, bc.Blocks, "empty block")
// 	assert.Equal(t, 2, len(bc.Blocks), "chain height not 2")
// }

// func TestVerifyBlockChain(t *testing.T) {
// 	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node2 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node3 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node4 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node5 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node6 := merkletree.NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node7 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
// 	node8 := merkletree.NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})

// 	tree1 := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node1, node2, node3, node4})
// 	tree2 := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node5, node6, node7, node8})

// 	blockchain := NewBlockChain()
// 	blockchain.AddBlock(*tree1)
// 	blockchain.AddBlock(*tree2)
// 	assert.True(t, blockchain.VerifyBlockChain(), "blockchain not verified")

// 	//modify a single value
// 	blockchain.Blocks[1].Nonce = 40
// 	assert.False(t, blockchain.VerifyBlockChain(), "blockchain verified")
// }
