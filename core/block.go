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
)

//! Errors
var ErrHashModification = errors.New("Attempt to modify hash value of block")

type Block struct {
	Timestamp     int64    `json:"timestamp"`
	Jobs          []byte   `json:"jobs"` //! would be modified based on job engine
	PrevBlockHash []byte   `json:"prevBlockHash"`
	MerkleRoot    []byte   `json:"merkleRoot"` //hash of merkle tree of jobs
	Hash          []byte   `json:"hash"`
	Difficulty    *big.Int `json:"difficulty"`
	Nonce         uint64   `json:"nonce"`
}

func NewBlock(jobs, pHash, mHash []byte) *Block {
	//! pow has to set nonce
	Block := &Block{
		Timestamp:     time.Now().Unix(),
		Jobs:          jobs,
		PrevBlockHash: pHash,
		MerkleRoot:    mHash,
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
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Jobs, timestamp, b.MerkleRoot, []byte(strconv.FormatInt(int64(b.Nonce), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	if reflect.ValueOf(b.Hash).IsNil() {
		b.Hash = hash[:]
		return nil
	}
	return ErrHashModification
}

func (b *Block) VerifyBlock() bool {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Jobs, timestamp, b.MerkleRoot, []byte(strconv.FormatInt(int64(b.Nonce), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	return string(hash[:]) == string(b.Hash)
}
