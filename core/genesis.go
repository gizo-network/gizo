package core

import (
	"encoding/hex"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/kpango/glg"
)

//! modify on job engine creation

//GenesisBlock returns genesis block
func GenesisBlock(by string) *Block {
	glg.Info("Core: Creating Genesis Block")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func Genesis(){return 1+1}", "Genesis", false, hex.EncodeToString(priv))
	node := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	tree := merkletree.MerkleTree{
		Root:      node.GetHash(),
		LeafNodes: []*merkletree.MerkleNode{node},
	}
	prevHash := []byte("00000000000000000000000000000000000000")
	block := NewBlock(tree, prevHash, 0, 10, by)
	return block
}
