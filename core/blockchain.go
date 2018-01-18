package core

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gizo-network/gizo/core/merkle_tree"

	"github.com/kpango/glg"

	"github.com/boltdb/bolt"
)

type BlockChain struct {
	tip []byte
	db  *bolt.DB
	mu  sync.RWMutex
}

func (bc BlockChain) GetTip() []byte {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.tip
}

func (bc *BlockChain) SetTip(t []byte) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.tip = t
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

func CreateBlockChain() *BlockChain {
	dbFile := path.Join(IndexPath, fmt.Sprintf(IndexDB, "testnodeid")) //FIXME: integrate node id
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
			Header:    genesis.GetHeader(),
			Height:    genesis.GetHeight(),
			TotalJobs: uint(len(genesis.GetJobs())),
			FileName:  genesis.FileStats().Name(),
			FileSize:  genesis.FileStats().Size(),
		}
		blockinfoBytes := blockinfo.Serialize()

		if err = b.Put(genesis.Header.GetHash(), blockinfoBytes); err != nil {
			glg.Fatal(err)
		}

		//latest block on the chain
		if err = b.Put([]byte("l"), genesis.Header.GetHash()); err != nil {
			glg.Fatal(err)
		}
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	bc := &BlockChain{
		tip: genesis.Header.GetHash(),
		db:  db,
	}
	return bc
}

func (bc *BlockChain) GetBlockInfo(hash []byte) *BlockInfo {
	var blockinfo *BlockInfo
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		blockinfoBytes := b.Get(hash)
		blockinfo = DeserializeBlockInfo(blockinfoBytes)
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	return blockinfo
}

func (bc *BlockChain) GetLatestHeight() uint64 {
	var lastBlock *BlockInfo
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		lastBlockBytes := b.Get(bc.GetTip())
		lastBlock = DeserializeBlockInfo(lastBlockBytes)
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	return lastBlock.GetHeight()
}

func (bc *BlockChain) AddBlock(block *Block) {
	if block.VerifyBlock() == false {
		glg.Warn("Unverified block cannot be added to the blockchain")
		return
	}
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		inDb := b.Get(block.Header.GetHash())
		if inDb != nil {
			glg.Warn("Block exists in blockchain")
			return nil
		}

		blockinfo := BlockInfo{
			Header:    block.GetHeader(),
			Height:    block.GetHeight(),
			TotalJobs: uint(len(block.GetJobs())),
			FileName:  block.FileStats().Name(),
			FileSize:  block.FileStats().Size(),
		}

		if err := b.Put(block.GetHeader().GetHash(), blockinfo.Serialize()); err != nil {
			glg.Fatal(err)
		}

		if block.GetHeight() > bc.GetBlockInfo(bc.GetTip()).GetHeight() {
			if err := b.Put([]byte("l"), block.GetHeader().GetHash()); err != nil {
				glg.Fatal(err)
			}
			bc.SetTip(block.GetHeader().GetHash())
		}
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
}

func (bc BlockChain) iterator() *BlockChainIterator {
	return &BlockChainIterator{
		current: bc.tip,
		db:      bc.db,
	}
}

func (bc *BlockChain) FindJob(h []byte) bool {
	var tree merkle_tree.MerkleTree
	bci := bc.iterator()
	for {
		block := bci.Next()
		if block.GetHeight() == 0 {
			return false
		}
		tree.SetLeafNodes(block.GetJobs())
		found, err := tree.Search(h)
		if err != nil {
			glg.Fatal(err)
		}
		if found {
			return true
		}
	}
}

func (bc *BlockChain) GetBlockHashes() [][]byte {
	var hashes [][]byte
	bci := bc.iterator()
	for {
		block := bci.Next()
		hashes = append(hashes, block.GetHeader().GetHash())
		if block.GetHeight() == 0 {
			break
		}
	}
	return hashes
}

func dbExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
