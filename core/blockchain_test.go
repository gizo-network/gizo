package core

import (
	"testing"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

//FIXME: Debug bc iterator

func TestNewBlockChain(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
	bc := CreateBlockChain()
	assert.NotNil(t, bc)
}

func TestAddBlock(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	err := bc.AddBlock(block)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(bc.GetLatestHeight()))
}

func TestVerify(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
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
	bc := CreateBlockChain()
	block := NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestBlock().GetHeight()+1, 10)
	bc.AddBlock(block)
	assert.True(t, bc.Verify())
}

func TestGetBlockInfo(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	bc.AddBlock(block)
	blockinfo, err := bc.GetBlockInfo(block.GetHeader().GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, blockinfo)
}

func TestGetBlocksWithinMinute(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
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
	bc := CreateBlockChain()
	block := NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestBlock().GetHeight()+1, 10)
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetBlocksWithinMinute())
}

func TestGetLatestHeight(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
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
	bc := CreateBlockChain()
	block := NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestBlock().GetHeight()+1, 10)
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetLatestHeight())
}

func TestFindJob(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	bc := CreateBlockChain()
	bc.AddBlock(block)
	f, err := bc.FindJob(node5.GetJob().GetID())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestGetBlockHashes(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
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
	bc := CreateBlockChain()
	block := NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestBlock().GetHeight()+1, 10)
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetBlockHashes())
}

func TestCreateBlockChain(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
	bc := CreateBlockChain()
	assert.NotNil(t, bc)
}
