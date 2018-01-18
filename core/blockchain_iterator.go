package core

import (
	"github.com/boltdb/bolt"
	"github.com/kpango/glg"
)

//pass db from blockchain
type BlockChainIterator struct {
	current []byte
	db      *bolt.DB
}

func (i *BlockChainIterator) SetCurrent(c []byte) {
	i.current = c
}

func (i BlockChainIterator) GetCurrent() []byte {
	return i.current
}

func (i *BlockChainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		blockinfoBytes := b.Get(i.GetCurrent())
		blockinfo := DeserializeBlockInfo(blockinfoBytes)
		block = blockinfo.GetBlock()
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	i.SetCurrent(block.GetHeader().GetPrevBlockHash())
	return block
}
