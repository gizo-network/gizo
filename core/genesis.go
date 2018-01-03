package core

import "github.com/gizo-network/gizo/core/merkle_tree"

//! modify on job engine creation

func GenesisBlock() *Block {
	node := merkle_tree.NewNode([]byte("Create genesis block"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	prevHash := []byte("00000000000000000000000000000000000000")
	block := NewBlock(*node, prevHash)
	return block
}
