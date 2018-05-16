package core

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/helpers"

	"github.com/kpango/glg"
)

var (
	ErrUnableToExport   = errors.New("Unable to export block")
	ErrHashModification = errors.New("Attempt to modify hash value of block")
)

//Block - the foundation of blockchain
type Block struct {
	Header     BlockHeader
	Jobs       []*merkletree.MerkleNode
	Height     uint64
	ReceivedAt int64  //time it was received
	By         string // id of node that generated block
}

//GetHeader returns the block header
func (b Block) GetHeader() BlockHeader {
	return b.Header
}

//sets the block header
func (b *Block) setHeader(h BlockHeader) {
	b.Header = h
}

//GetNodes retuns the block's merklenodes
func (b Block) GetNodes() []*merkletree.MerkleNode {
	return b.Jobs
}

//sets merklenodes
func (b *Block) setNodes(j []*merkletree.MerkleNode) {
	b.Jobs = j
}

//GetHeight returns the block height
func (b Block) GetHeight() uint64 {
	return b.Height
}

//sets the block height
func (b *Block) setHeight(h uint64) {
	b.Height = h
}

//NewBlock returns a new block
func NewBlock(tree merkletree.MerkleTree, pHash []byte, height uint64, difficulty uint8, by string) *Block {
	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: pHash,
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(difficulty)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: height,
		By:     by,
	}
	pow := NewPOW(block)
	pow.run() //! mines block
	err := block.Export()
	if err != nil {
		glg.Fatal(err)
	}
	return block
}

//writes block on disk
func (b Block) Export() error {
	glg.Info("Core: Exporting block - " + hex.EncodeToString(b.GetHeader().GetHash()))
	InitializeDataPath()
	if b.IsEmpty() {
		return ErrUnableToExport
	}
	bBytes := b.Serialize()
	var err error
	if os.Getenv("ENV") == "dev" {
		err = ioutil.WriteFile(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, hex.EncodeToString(b.Header.GetHash()))), []byte(helpers.Encode64(bBytes)), os.FileMode(0555))
	} else {
		err = ioutil.WriteFile(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, hex.EncodeToString(b.Header.GetHash()))), []byte(helpers.Encode64(bBytes)), os.FileMode(0555))
	}
	if err != nil {
		glg.Fatal(err)
	}
	return nil
}

//Import reads block file into memory
func (b *Block) Import(hash []byte) {
	glg.Info("Core: Importing block - " + hex.EncodeToString(hash))
	if b.IsEmpty() == false {
		glg.Warn("Overwriting umempty block")
	}
	var read []byte
	var err error
	if os.Getenv("ENV") == "dev" {
		read, err = ioutil.ReadFile(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, hex.EncodeToString(hash))))
	} else {
		read, err = ioutil.ReadFile(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, hex.EncodeToString(hash))))
	}
	if err != nil {
		glg.Fatal(err) //FIXME: handle block doesn't exist by asking peer
	}
	bBytes := helpers.Decode64(string(read))
	temp, err := DeserializeBlock(bBytes)
	if err != nil {
		glg.Fatal(err)
	}
	b.setHeader(temp.GetHeader())
	b.setHeight(temp.GetHeight())
	b.setNodes(temp.GetNodes())
}

//returns the file stats of a blockfile
func (b Block) fileStats() os.FileInfo {
	var info os.FileInfo
	var err error
	if os.Getenv("ENV") == "dev" {
		info, err = os.Stat(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, hex.EncodeToString(b.Header.GetHash()))))
	} else {
		info, err = os.Stat(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, hex.EncodeToString(b.Header.GetHash()))))
	}
	if os.IsNotExist(err) {
		glg.Fatal("Block file doesn't exist")
	}
	return info
}

//IsEmpty returns true is block is empty
func (b *Block) IsEmpty() bool {
	return b.Header.GetHash() == nil
}

//Serialize returns bytes of block
func (b *Block) Serialize() []byte {
	temp, err := json.Marshal(*b)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}

//DeserializeBlock returns block from bytes
func DeserializeBlock(b []byte) (*Block, error) {
	var temp Block
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return nil, err
	}
	return &temp, nil
}

//VerifyBlock verifies a block
func (b *Block) VerifyBlock() bool {
	glg.Info("Core: Verifying block - " + hex.EncodeToString(b.GetHeader().GetHash()))
	pow := NewPOW(b)
	return pow.Validate()
}

//DeleteFile deletes block file on disk
func (b Block) DeleteFile() {
	glg.Info("Core: Deleting blockfile - " + hex.EncodeToString(b.GetHeader().GetHash()))
	var err error
	if os.Getenv("ENV") == "dev" {
		err = os.Remove(path.Join(BlockPathDev, b.fileStats().Name()))
	} else {
		err = os.Remove(path.Join(BlockPathProd, b.fileStats().Name()))
	}
	if err != nil {
		glg.Fatal(err)
	}
}
