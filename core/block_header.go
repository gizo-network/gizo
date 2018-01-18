package core

import "math/big"

type BlockHeader struct {
	Timestamp     int64    `json:"timestamp"`
	PrevBlockHash []byte   `json:"prevBlockHash"`
	MerkleRoot    []byte   `json:"merkleroot"`
	Nonce         uint64   `json:"nonce"`
	Difficulty    *big.Int `json:"difficulty"`
	Hash          []byte   `json:"hash"`
}

func (bh BlockHeader) GetTimestamp() int64 {
	return bh.Timestamp
}

func (bh *BlockHeader) SetTimestamp(t int64) {
	bh.Timestamp = t
}

func (bh BlockHeader) GetPrevBlockHash() []byte {
	return bh.PrevBlockHash
}

func (bh *BlockHeader) SetPrevBlockHash(h []byte) {
	bh.PrevBlockHash = h
}

func (bh BlockHeader) GetMerkleRoot() []byte {
	return bh.MerkleRoot
}

func (bh *BlockHeader) SetMerkleRoot(mr []byte) {
	bh.MerkleRoot = mr
}

func (bh BlockHeader) GetNonce() uint64 {
	return bh.Nonce
}

func (bh *BlockHeader) SetNonce(n uint64) {
	bh.Nonce = n
}

func (bh BlockHeader) GetDifficulty() big.Int {
	return *bh.Difficulty
}

func (bh *BlockHeader) SetDifficulty(d big.Int) {
	bh.Difficulty = &d
}

func (bh BlockHeader) GetHash() []byte {
	return bh.Hash
}

func (bh *BlockHeader) SetHash(h []byte) {
	bh.Hash = h
}
