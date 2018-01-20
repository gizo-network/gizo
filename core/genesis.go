package core

import "github.com/gizo-network/gizo/core/merkletree"

//! modify on job engine creation

func GenesisBlock() *Block {
	node := merkletree.NewNode([]byte("Create genesis block"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	tree := merkletree.MerkleTree{
		Root:      node.GetHash(),
		LeafNodes: []*merkletree.MerkleNode{node},
	}
	prevHash := []byte("00000000000000000000000000000000000000")
	block := NewBlock(tree, prevHash, 0, 5)
	return block
}
