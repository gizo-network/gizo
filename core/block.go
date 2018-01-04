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
	MerkleRoot    merkle_tree.MerkleNode `json:"merkleRoot"`
	Hash          []byte                 `json:"hash"`
	Difficulty    *big.Int               `json:"difficulty"`
	Nonce         uint64                 `json:"nonce"`
}

func NewBlock(mRoot merkle_tree.MerkleNode, pHash []byte) *Block {
	//! pow has to set nonce
	//! dificullty engine would set difficulty
	// Log.Logger.Info("Creating new block")
	Block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: pHash,
		MerkleRoot:    mRoot,
	}
	err := Block.SetHash()
	if err != nil {
		glg.Fatal(err)
	}
	return Block
}

func MarshalBlock(b *Block) ([]byte, error) {
	temp, err := json.Marshal(*b)
	if err != nil {
		return nil, err
	}
	return temp, nil
}

func UnMashalBlock(b []byte) (*Block, error) {
	var temp Block
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return nil, err
	}
	return &temp, nil
}

func (b *Block) SetHash() error {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	mBytes, err := merkle_tree.MarshalMerkleNode(b.MerkleRoot)
	if err != nil {
		glg.Fatal(err)
	}
	headers := bytes.Join([][]byte{b.PrevBlockHash, timestamp, mBytes, []byte(strconv.FormatInt(int64(b.Nonce), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	if reflect.ValueOf(b.Hash).IsNil() {
		b.Hash = hash[:]
		return nil
	}
	return ErrHashModification
}

func (b *Block) VerifyBlock() bool {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	mBytes, err := merkle_tree.MarshalMerkleNode(b.MerkleRoot)
	if err != nil {
		glg.Fatal(err)
	}
	headers := bytes.Join([][]byte{b.PrevBlockHash, timestamp, mBytes, []byte(strconv.FormatInt(int64(b.Nonce), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	return string(hash[:]) == string(b.Hash)
}
