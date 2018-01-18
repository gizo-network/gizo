package core

import "github.com/gizo-network/gizo/core/merkle_tree"

//! modify on job engine creation

func GenesisBlock() *Block {
	node := merkle_tree.NewNode([]byte("Create genesis block"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	tree := merkle_tree.MerkleTree{
		Root:      node.GetHash(),
		LeafNodes: []*merkle_tree.MerkleNode{node},
	}
	prevHash := []byte("00000000000000000000000000000000000000")
	block := NewBlock(tree, prevHash, 0)
	return block
}
