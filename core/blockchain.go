package core

import (
	"github.com/boltdb/bolt"
)

var (
	VALID_TREE  = 1
	VALID_BLOCK = 2
)

type BlockChain struct {
	tip Block
	db  *bolt.DB
}

type BlockInfo struct {
	Header     BlockHeader `json:"header"`
	Height     uint        `json:"height"`
	TotalJobs  uint        `json:"total_jobs"`
	Validation uint        `json:"validation"`
	FileName   string      `json:"file_name"`
	FileSize   int64       `json:"file_size"`
}

// func (bc *BlockChain) AddBlock(tree merkle_tree.MerkleTree) {
// 	//FIXME: get current height from db and add 1
// 	prevBlock := bc.Blocks[len(bc.Blocks)-1]
// 	newBlock := NewBlock(tree, prevBlock.Hash, 0)
// 	bc.Blocks = append(bc.Blocks, newBlock)
// }

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

// func NewBlockChain() *BlockChain {
// 	dbFile := path.Join(BlockPath, fmt.Sprintf(BlockFile, "6546546351as3dfasdfas6d"))
// 	bc := &BlockChain{}
// 	bc.Blocks = append(bc.Blocks, GenesisBlock())
// 	return bc
// }
