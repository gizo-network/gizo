package core

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"strconv"

	"github.com/gizo-network/gizo/core/consensus"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/kpango/glg"
)

var maxNonce = math.MaxInt64

type POW struct {
	block  *Block
	target *big.Int
}

func (p *POW) SetBlock(b *Block) {
	p.block = b
}

func (p POW) GetBlock() *Block {
	return p.block
}

func (p *POW) SetTarget(t *big.Int) {
	p.target = t
}

func (p POW) GetTarget() *big.Int {
	return p.target
}

func (p POW) prepareData(nonce int) []byte {
	tree := merkletree.MerkleTree{Root: p.block.GetHeader().GetMerkleRoot(), LeafNodes: p.block.GetJobs()}
	mBytes, err := tree.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	data := bytes.Join(
		[][]byte{
			p.block.GetHeader().GetPrevBlockHash(),
			[]byte(strconv.FormatInt(p.block.Header.GetTimestamp(), 10)),
			mBytes,
			[]byte(strconv.FormatInt(int64(nonce), 10)),
			[]byte(strconv.FormatInt(int64(p.block.GetHeight()), 10)),
			[]byte(strconv.FormatInt(int64(consensus.Difficulty), 10)),
		},
		[]byte{},
	)
	return data
}

func (p *POW) Run() {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	glg.Info("Mining block")
	for nonce < maxNonce {
		glg.Info(nonce)
		hash = sha256.Sum256(p.prepareData(nonce))
		glg.Infof("%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	p.block.Header.SetDifficulty(*big.NewInt(int64(consensus.Difficulty)))
	p.block.Header.SetHash(hash[:])
	p.block.Header.SetNonce(uint64(nonce))
}

func NewPOW(b *Block) *POW {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-consensus.Difficulty))

	pow := &POW{
		target: target,
		block:  b,
	}
	return pow
}
