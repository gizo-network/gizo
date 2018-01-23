package core

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gizo-network/gizo/core/merkletree"

	"github.com/kpango/glg"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/now"
)

type BlockChain struct {
	tip []byte
	db  *bolt.DB
	mu  *sync.RWMutex
}

func (bc *BlockChain) GetTip() []byte {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.tip
}

func (bc *BlockChain) SetTip(t []byte) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.tip = t
}

func (bc *BlockChain) DB() *bolt.DB {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.db
}

func (bc *BlockChain) SetDB(db *bolt.DB) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.db = db
}

//GetBlockInfo returns the blockinfo of a particular block from the db
func (bc *BlockChain) GetBlockInfo(hash []byte) *BlockInfo {
	var blockinfo *BlockInfo
	err := bc.DB().View(func(tx *bolt.Tx) error {
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

//GetBlocksWithinMinute returns all blocks in the db within the last minute
func (bc *BlockChain) GetBlocksWithinMinute() []Block {
	var blocks []Block
	now := now.New(time.Now())

	bci := bc.iterator()
	for {
		block := bci.Next()
		if block.GetHeight() == 0 && block.GetHeader().GetTimestamp() > now.BeginningOfMinute().Unix() {
			blocks = append(blocks, *block)
			break
		} else if block.GetHeader().GetTimestamp() > now.BeginningOfMinute().Unix() {
			blocks = append(blocks, *block)
		} else {
			break
		}
	}
	return blocks
}

//GetLatestHeight returns the height of the latest block to the blockchain
func (bc *BlockChain) GetLatestHeight() uint64 {
	var lastBlock *BlockInfo
	err := bc.DB().View(func(tx *bolt.Tx) error {
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

//AddBlock adds block to the blockchain
func (bc *BlockChain) AddBlock(block *Block) {
	if block.VerifyBlock() == false {
		glg.Warn("Unverified block cannot be added to the blockchain")
		return
	}
	err := bc.DB().Update(func(tx *bolt.Tx) error {
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

		//FIXME: handle a fork
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

// return a BlockChainIterator to loop throught the blockchain
func (bc *BlockChain) iterator() *BlockChainIterator {
	return &BlockChainIterator{
		current: bc.GetTip(),
		db:      bc.DB(),
	}
}

//FindJob returns the merklenode of a job from the blockchain
func (bc *BlockChain) FindJob(h []byte) *merkletree.MerkleNode {
	var tree merkletree.MerkleTree
	bci := bc.iterator()
	for {
		block := bci.Next()
		if block.GetHeight() == 0 {
			return nil
		}
		tree.SetLeafNodes(block.GetJobs())
		found, err := tree.Search(h)
		if err != nil {
			glg.Fatal(err)
		}
		return found
	}
}

//GetBlockHashes returns all the hashes of all the blocks in the current bc
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

//CreateBlockChain initializes a db, set's the tip to GenesisBlock and returns the blockchain
func CreateBlockChain() *BlockChain {
	InitializeDataPath()
	dbFile := path.Join(IndexPath, fmt.Sprintf(IndexDB, "testnodeid")) //FIXME: integrate node id
	if dbExists(dbFile) {
		var tip []byte
		glg.Warn("Using existing blockchain")
		db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
		if err != nil {
			glg.Fatal(err)
		}
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BlockBucket))
			tip = b.Get([]byte("l"))
			return nil
		})
		if err != nil {
			glg.Fatal(err)
		}
		return &BlockChain{
			tip: tip,
			db:  db,
			mu:  &sync.RWMutex{},
		}
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
		mu:  &sync.RWMutex{},
	}
	return bc
}

func dbExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
