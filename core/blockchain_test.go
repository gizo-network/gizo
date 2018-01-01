package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockChain(t *testing.T) {
	bc := NewBlockChain()
	assert.NotNil(t, bc)
	assert.NotNil(t, bc.Blocks)
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockChain()
	bc.AddBlock([]byte("test block"), []byte("merklehash"))
	assert.NotEmpty(t, bc.Blocks, "empty block")
	assert.Equal(t, 2, len(bc.Blocks), "chain height not 2")
}

func TestVerifyBlockChain(t *testing.T) {
	blockchain := NewBlockChain()
	blockchain.AddBlock([]byte("jobs example 1"), []byte("merklehash"))
	blockchain.AddBlock([]byte("jobs example 2"), []byte("merklehash"))
	assert.True(t, blockchain.VerifyBlockChain(), "blockchain not verified")

	//modify a single value
	blockchain.Blocks[1].Nonce = 40
	assert.False(t, blockchain.VerifyBlockChain(), "blockchain verified")
}
