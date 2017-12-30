package core

import (
	"bytes"
	"crypto/sha256"
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
	MerkleHash    []byte   `json:"merkleHash"` //hash of merkle tree of jobs
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
		MerkleHash:    mHash,
	}
	return Block
}

func (b *Block) SetHash() error {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Jobs, timestamp, b.MerkleHash, []byte(strconv.FormatInt(int64(b.Nonce), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	if reflect.ValueOf(b.Hash).IsNil() {
		b.Hash = hash[:]
		return nil
	}
	return ErrHashModification
}

func (b *Block) VerifyBlock() bool {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Jobs, timestamp, b.MerkleHash, []byte(strconv.FormatInt(int64(b.Nonce), 10))}, []byte{})
	hash := sha256.Sum256(headers)
	return string(hash[:]) == string(b.Hash)
}
