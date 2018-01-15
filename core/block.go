package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core/merkle_tree"
)

var (
	ErrUnableToExport   = errors.New("Unable to export block")
	ErrHashModification = errors.New("Attempt to modify hash value of block")
)

// var Log = helpers.NewLogger()

type Block struct {
	Header BlockHeader               `json:"header"`
	Jobs   []*merkle_tree.MerkleNode `json:"jobs"`
	Height uint64                    `json:"height"`
}

type BlockHeader struct {
	Timestamp     int64    `json:"timestamp"`
	PrevBlockHash []byte   `json:"prevBlockHash"`
	MerkleRoot    []byte   `json:"merkleroot"`
	Nonce         uint64   `json:"nonce"`
	Difficulty    *big.Int `json:"difficulty"`
	Hash          []byte   `json:"hash"`
}

//FIXME: implement block status
func NewBlock(tree merkle_tree.MerkleTree, pHash []byte, height uint64) *Block {
	//! pow has to set nonce
	//! dificullty engine would set difficulty
	// Log.Logger.Info("Creating new block")
	Block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: pHash,
			MerkleRoot:    tree.Root,
		},
		Jobs:   tree.LeafNodes,
		Height: height,
	}
	err := Block.setHash()
	if err != nil {
		glg.Fatal(err)
	}
	return Block
}

func (b *Block) Export() error {
	if b.IsEmpty() {
		return ErrUnableToExport
	}
	bBytes, err := b.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf(BlockFile, b.Header.Hash), []byte(helpers.Encode64(bBytes)), os.FileMode(0555))
	if err != nil {
		glg.Fatal(err)
	}
	return nil
}

func (b *Block) IsEmpty() bool {
	return reflect.ValueOf(b.Jobs).IsNil() == true && reflect.ValueOf(b.Height).IsNil() == true && reflect.ValueOf(b.Header).IsNil() == true
}

//Serialize returns bytes of block
func (b *Block) Serialize() ([]byte, error) {
	temp, err := json.Marshal(*b)
	if err != nil {
		return nil, err
	}
	return temp, nil
}

//DeserializeBlock returns block of bytes
func DeserializeBlock(b []byte) (*Block, error) {
	var temp Block
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return nil, err
	}
	return &temp, nil
}

func (b *Block) setHash() error {
	timestamp := []byte(strconv.FormatInt(b.Header.Timestamp, 10))
	tree := merkle_tree.MerkleTree{Root: b.Header.MerkleRoot, LeafNodes: b.Jobs}
	mBytes, err := tree.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	headers := bytes.Join([][]byte{b.Header.PrevBlockHash, timestamp, mBytes, []byte(strconv.FormatInt(int64(b.Header.Nonce), 10)), []byte(strconv.FormatInt(int64(b.Height), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	if reflect.ValueOf(b.Header.Hash).IsNil() {
		b.Header.Hash = hash[:]
		return nil
	}
	return ErrHashModification
}

func (b *Block) VerifyBlock() bool {
	timestamp := []byte(strconv.FormatInt(b.Header.Timestamp, 10))
	tree := merkle_tree.MerkleTree{Root: b.Header.MerkleRoot, LeafNodes: b.Jobs}
	mBytes, err := tree.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	headers := bytes.Join([][]byte{b.Header.PrevBlockHash, timestamp, mBytes, []byte(strconv.FormatInt(int64(b.Header.Nonce), 10)), []byte(strconv.FormatInt(int64(b.Height), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	return bytes.Equal(hash[:], b.Header.Hash)
}
