package core

import (
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
)

//! modify on job engine creation

func GenesisBlock() *Block {
	j := job.NewJob("func test(){return 1+1}", "test")
	node := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	tree := merkletree.MerkleTree{
		Root:      node.GetHash(),
		LeafNodes: []*merkletree.MerkleNode{node},
	}
	prevHash := []byte("00000000000000000000000000000000000000")
	block := NewBlock(tree, prevHash, 0, 10)
	return block
}
