package core

import (
	"github.com/gizo-network/gizo/core/merkle_tree"
)

type BlockChain struct {
	Blocks []*Block `json:"blocks"`
}

func (bc *BlockChain) AddBlock(tree merkle_tree.MerkleTree) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(*tree.Root, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// VerifyBlockChain returns true if blockchain is valid
func (bc *BlockChain) VerifyBlockChain() bool {
	for i, val := range bc.Blocks {
		if i != 0 {
			if val.VerifyBlock() == false || (string(val.PrevBlockHash) == string(bc.Blocks[i-1].Hash)) == false {
				return false
			}
		}
	}
	return true
}

func NewBlockChain() *BlockChain {
	bc := &BlockChain{}
	bc.Blocks = append(bc.Blocks, GenesisBlock())
	return bc
}
