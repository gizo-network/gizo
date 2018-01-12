package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"math/big"
	"reflect"
	"strconv"
	"time"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core/merkle_tree"
)

var ErrHashModification = errors.New("Attempt to modify hash value of block")

// var Log = helpers.NewLogger()

type Block struct {
	Timestamp     int64                  `json:"timestamp"`
	PrevBlockHash []byte                 `json:"prevBlockHash"`
	MerkleTree    merkle_tree.MerkleTree `json:"merkletree"`
	Hash          []byte                 `json:"hash"`
	Difficulty    *big.Int               `json:"difficulty"`
	Height        uint64                 `json:"height"`
	Nonce         uint64                 `json:"nonce"`
}

func NewBlock(tree merkle_tree.MerkleTree, pHash []byte, height uint64) *Block {
	//! pow has to set nonce
	//! dificullty engine would set difficulty
	// Log.Logger.Info("Creating new block")
	Block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: pHash,
		MerkleTree:    tree,
		Height:        height,
	}
	err := Block.setHash()
	if err != nil {
		glg.Fatal(err)
	}
	return Block
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
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	mBytes, err := b.MerkleTree.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	headers := bytes.Join([][]byte{b.PrevBlockHash, timestamp, mBytes, []byte(strconv.FormatInt(int64(b.Nonce), 10)), []byte(strconv.FormatInt(int64(b.Height), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	if reflect.ValueOf(b.Hash).IsNil() {
		b.Hash = hash[:]
		return nil
	}
	return ErrHashModification
}

func (b *Block) VerifyBlock() bool {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	mBytes, err := b.MerkleTree.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	headers := bytes.Join([][]byte{b.PrevBlockHash, timestamp, mBytes, []byte(strconv.FormatInt(int64(b.Nonce), 10)), []byte(strconv.FormatInt(int64(b.Height), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	return bytes.Equal(hash[:], b.Hash)
}
