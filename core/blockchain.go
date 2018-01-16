package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/kpango/glg"

	"github.com/boltdb/bolt"
)

// var (
// 	VALID_TREE  = 1
// 	VALID_BLOCK = 2
// )

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

type BlockInfo struct {
	Header    BlockHeader `json:"header"`
	Height    uint64      `json:"height"`
	TotalJobs uint        `json:"total_jobs"`
	// Validation uint        `json:"validation"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

func (b *BlockInfo) Serialize() ([]byte, error) {
	temp, err := json.Marshal(*b)
	if err != nil {
		return nil, err
	}
	return temp, nil
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

func CreateBlockChain() *BlockChain {
	dbFile := path.Join(IndexPath, fmt.Sprintf(IndexPath, "testnodeid")) //FIXME: integrate node id
	if dbExists(dbFile) {
		glg.Fatal("Blockchain exists")
	}
	genesis := GenesisBlock()
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
	if err != nil {
		glg.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BlockBucket))
		if err != nil {
			glg.Fatal(err)
		}
		blockinfo := BlockInfo{
			Header:    genesis.Header,
			Height:    genesis.Height,
			TotalJobs: uint(len(genesis.Jobs)),
			FileName:  genesis.FileStats().Name(),
			FileSize:  genesis.FileStats().Size(),
		}
		blockinfoBytes, err := blockinfo.Serialize()
		if err != nil {
			glg.Fatal(err)
		}
		err = b.Put(genesis.Header.Hash, blockinfoBytes)
		if err != nil {
			glg.Fatal(err)
		}
		err = b.Put([]byte("l"), genesis.Header.Hash) //latest block on the chain
		if err != nil {
			glg.Fatal(err)
		}
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	bc := &BlockChain{
		tip: genesis.Header.Hash,
		db:  db,
	}
	return bc
}

func (bc *BlockChain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockBucket))
		inDb := bucket.Get(block.Header.Hash)
		if inDb == nil {
			glg.Warn("Block is already in blockchain")
		}

		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
}

func dbExists(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *BlockChain) GetTip() []byte {
	return bc.tip
}
