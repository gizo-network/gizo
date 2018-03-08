package core

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gizo-network/gizo/job"

	"github.com/gizo-network/gizo/core/merkletree"

	"github.com/kpango/glg"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/now"
)

var (
	ErrUnverifiedBlock = errors.New("Unverified block cannot be added to the blockchain")
	ErrJobNotFound     = errors.New("Job not found")
	ErrBlockNotFound   = errors.New("Blockinfo not found")
)

//BlockChain - singly linked list of blocks
type BlockChain struct {
	tip []byte //! hash of latest block in the blockchain
	db  *bolt.DB
	mu  *sync.RWMutex
}

//returns the blockinfo of the latest block in the blockchain
func (bc *BlockChain) getTip() []byte {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.tip
}

//sets the tip
func (bc *BlockChain) setTip(t []byte) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.tip = t
}

//returns the db
func (bc *BlockChain) getDB() *bolt.DB {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.db
}

//sets the db
func (bc *BlockChain) setDB(db *bolt.DB) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.db = db
}

//GetBlockInfo returns the blockinfo of a particular block from the db
func (bc *BlockChain) GetBlockInfo(hash []byte) (*BlockInfo, error) {
	glg.Info("Core: Getting blockinfo - " + hex.EncodeToString(hash))
	var blockinfo *BlockInfo
	err := bc.getDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		blockinfoBytes := b.Get(hash)
		if blockinfoBytes != nil {
			blockinfo = DeserializeBlockInfo(blockinfoBytes)
		} else {
			blockinfo = nil
		}
		return nil
	})
	if err != nil {
		glg.Fatal(err) //! handle db failure error
	}
	if blockinfo != nil {
		return blockinfo, nil
	}
	return nil, ErrBlockNotFound
}

func (bc BlockChain) GetPrevHash() []byte {
	return bc.GetLatestBlock().GetHeader().GetHash()
}

//GetBlocksWithinMinute returns all blocks in the db within the last minute
func (bc *BlockChain) GetBlocksWithinMinute() []Block {
	glg.Info("Core: Getting blocks within last minute")
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

//GetLatest15 retuns the latest 15 blocks
func (bc *BlockChain) GetLatest15() []Block {
	glg.Info("Core: Getting blocks within last minute")
	var blocks []Block
	bci := bc.iterator()
	for {
		if len(blocks) <= 15 {
			block := bci.Next()
			if block.GetHeight() == 0 {
				blocks = append(blocks, *block)
				break
			} else {
				blocks = append(blocks, *block)
			}
		} else {
			break
		}
	}
	return blocks
}

//GetLatestHeight returns the height of the latest block to the blockchain
func (bc *BlockChain) GetLatestHeight() uint64 {
	glg.Info("Core: Getting latest block height")
	var lastBlock *BlockInfo
	err := bc.getDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		lastBlockBytes := b.Get(bc.getTip())
		lastBlock = DeserializeBlockInfo(lastBlockBytes)
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	return lastBlock.GetHeight()
}

//GetLatestBlock returns the tip as a block
func (bc *BlockChain) GetLatestBlock() *Block {
	glg.Info("Core: Getting latest block")
	var lastBlock *BlockInfo
	err := bc.getDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		lastBlockBytes := b.Get(bc.getTip())
		lastBlock = DeserializeBlockInfo(lastBlockBytes)
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	return lastBlock.GetBlock()
}

//GetNextHeight returns the next height in the blockchain
func (bc BlockChain) GetNextHeight() uint64 {
	return bc.GetLatestBlock().GetHeight() + 1
}

//AddBlock adds block to the blockchain
func (bc *BlockChain) AddBlock(block *Block) error {
	glg.Info("Core: Adding block to the blockchain - " + hex.EncodeToString(block.GetHeader().GetHash()))
	if block.VerifyBlock() == false {
		return ErrUnverifiedBlock
	}
	err := bc.getDB().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		inDb := b.Get(block.Header.GetHash())
		if inDb != nil {
			glg.Warn("Block exists in blockchain")
			return nil
		}

		blockinfo := BlockInfo{
			Header:    block.GetHeader(),
			Height:    block.GetHeight(),
			TotalJobs: uint(len(block.GetNodes())),
			FileName:  block.fileStats().Name(),
			FileSize:  block.fileStats().Size(),
		}

		if err := b.Put(block.GetHeader().GetHash(), blockinfo.Serialize()); err != nil {
			glg.Fatal(err)
		}

		//FIXME: handle a fork
		latest, err := bc.GetBlockInfo(bc.getTip())
		if err != nil {
			glg.Fatal(err)
		}
		if block.GetHeight() > latest.GetHeight() {
			if err := b.Put([]byte("l"), block.GetHeader().GetHash()); err != nil {
				glg.Fatal(err)
			}
			bc.setTip(block.GetHeader().GetHash())
		}
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	return nil
}

// return a BlockChainIterator to loop throught the blockchain
func (bc *BlockChain) iterator() *BlockChainIterator {
	return &BlockChainIterator{
		current: bc.getTip(),
		db:      bc.getDB(),
	}
}

//FindJob returns the job from the blockchain
func (bc *BlockChain) FindJob(id string) (*job.Job, error) {
	glg.Info("Core: Finding Job in the blockchain - " + id)
	var tree merkletree.MerkleTree
	bci := bc.iterator()
	for {
		block := bci.Next()
		if block.GetHeight() == 0 {
			return nil, ErrJobNotFound
		}
		tree.SetLeafNodes(block.GetNodes())
		found, err := tree.SearchJob(id)
		if err != nil {
			glg.Fatal(err)
		}
		return found, nil
	}
}

//FindMerkleNode returns the merklenode from the blockchain
func (bc *BlockChain) FindMerkleNode(h []byte) (*merkletree.MerkleNode, error) {
	glg.Info("Core: Finding merklenode - " + hex.EncodeToString(h))
	var tree merkletree.MerkleTree
	bci := bc.iterator()
	for {
		block := bci.Next()
		if block.GetHeight() == 0 {
			return nil, ErrJobNotFound
		}
		tree.SetLeafNodes(block.GetNodes())
		found, err := tree.SearchNode(h)
		if err != nil {
			glg.Fatal(err)
		}
		return found, nil
	}
}

//Verify verifies the blockchain
func (bc *BlockChain) Verify() bool {
	glg.Info("Core: Verifying Blockchain")
	bci := bc.iterator()
	for {
		block := bci.Next()
		if block.GetHeight() == 0 {
			return true
		}
		if block.VerifyBlock() == false {
			return false
		}
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
	glg.Info("Core: Creating blockchain database")
	InitializeDataPath()
	var dbFile string
	if os.Getenv("ENV") == "dev" {
		dbFile = path.Join(IndexPathDev, fmt.Sprintf(IndexDB, "testnodeid")) //FIXME: integrate node id
	} else {
		dbFile = path.Join(IndexPathProd, fmt.Sprintf(IndexDB, "testnodeid")) //FIXME: integrate node id
	}
	if dbExists(dbFile) {
		var tip []byte
		glg.Warn("Core: Using existing blockchain")
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
			TotalJobs: uint(len(genesis.GetNodes())),
			FileName:  genesis.fileStats().Name(),
			FileSize:  genesis.fileStats().Size(),
		}
		blockinfoBytes := blockinfo.Serialize()

		if err = b.Put(genesis.GetHeader().GetHash(), blockinfoBytes); err != nil {
			glg.Fatal(err)
		}

		//latest block on the chain
		if err = b.Put([]byte("l"), genesis.GetHeader().GetHash()); err != nil {
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
