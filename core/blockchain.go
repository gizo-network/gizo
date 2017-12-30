package core

type BlockChain struct {
	blocks []*Block `json:"blocks"`
}

func (bc *BlockChain) AddBlock(data string, jobs []byte) {
	//TODO: generate merkle tree hash from jobs
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(jobs, prevBlock.Hash, []byte("merkle hash temp"))
	bc.blocks = append(bc.blocks, newBlock)
}

func (bc *BlockChain) VerifyBlockChain() bool {
	for i, val := range bc.blocks {
		if i != 0 {
			if val.VerifyBlock() == false || (string(val.PrevBlockHash) == string(bc.blocks[i-1].Hash)) == false {
				return false
			}
		}
	}
	return true
}

func NewBlockChain() *BlockChain {
	bc := &BlockChain{}
	bc.blocks = append(bc.blocks, GenesisBlock)
	return bc
}
