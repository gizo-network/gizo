package core

type BlockChain struct {
	Blocks []*Block `json:"blocks"`
}

func (bc *BlockChain) AddBlock(jobs, merkleHash []byte) {
	//TODO: generate merkle tree hash from jobs
	//!FIXME: remove merklehash from arguments
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(jobs, prevBlock.Hash, merkleHash)
	newBlock.SetHash()
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
	GenesisBlock.SetHash()
	bc := &BlockChain{}
	bc.Blocks = append(bc.Blocks, GenesisBlock)
	return bc
}
