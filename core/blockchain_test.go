package core

import (
	"testing"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/stretchr/testify/assert"
)

//FIXME: Debug bc iterator

func TestNewBlockChain(t *testing.T) {
	RemoveDataPath()
	bc := CreateBlockChain()
	assert.NotNil(t, bc)
}

func TestAddBlock(t *testing.T) {
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	err := bc.AddBlock(block)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(bc.GetLatestHeight()))
}

func TestVerify(t *testing.T) {
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	bc.AddBlock(block)
	assert.True(t, bc.Verify())
}

func TestGetBlockInfo(t *testing.T) {
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	bc.AddBlock(block)
	blockinfo, err := bc.GetBlockInfo(block.GetHeader().GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, blockinfo)
}

func TestGetBlocksWithinMinute(t *testing.T) {
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetBlocksWithinMinute())
}

func TestGetLatestHeight(t *testing.T) {
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetLatestHeight())
}

func TestFindJob(t *testing.T) {
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	bc.AddBlock(block)
	j, err := bc.FindJob(node5.GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, j)
}

func TestGetBlockHashes(t *testing.T) {
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetBlockHashes())
}

func TestCreateBlockChain(t *testing.T) {
	RemoveDataPath()
	bc := CreateBlockChain()
	assert.NotNil(t, bc)
}
