package core_test

import (
	"testing"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestGetBlock(t *testing.T) {
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
	block := core.NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10)
	blockinfo := core.BlockInfo{
		Header:    block.GetHeader(),
		Height:    block.GetHeight(),
		TotalJobs: uint(len(block.GetNodes())),
		FileName:  block.FileStats().Name(),
		FileSize:  block.FileStats().Size(),
	}
	assert.Equal(t, block, blockinfo.GetBlock())
	block.DeleteFile()
}
