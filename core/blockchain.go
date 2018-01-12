package core

import (
	"github.com/boltdb/bolt"
	"github.com/gizo-network/gizo/core/merkle_tree"
)

type BlockChain struct {
	Blocks []*Block `json:"blocks"`
	db     *bolt.DB
}

func (bc *BlockChain) AddBlock(tree merkle_tree.MerkleTree) {
	//FIXME: get current height from db and add 1
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(tree, prevBlock.Hash, 0)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// VerifyBlockChain returns true if blockchain is valid
// func (bc *BlockChain) VerifyBlockChain() bool {
// 	for i, val := range bc.Blocks {
// 		if i != 0 {
// 			if val.VerifyBlock() == false || (string(val.PrevBlockHash) == string(bc.Blocks[i-1].Hash)) == false {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

func NewBlockChain() *BlockChain {
	// dbFile := path.Join(BlockPath, fmt.Sprintf(BlockFile, "6546546351as3dfasdfas6d"))
	bc := &BlockChain{}
	bc.Blocks = append(bc.Blocks, GenesisBlock())
	return bc
}
