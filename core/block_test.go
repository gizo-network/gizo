package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type String string

func TestNewBlock(t *testing.T) {
	jobs := []byte("test jobs")
	prevHash := []byte("00000000000000000000000000000000000000")
	mHash := []byte("0000000000000000000000000000000000000000")
	testBlock := NewBlock(jobs, prevHash, mHash)

	assert.NotNil(t, testBlock, "returned empty tblock")
	assert.Equal(t, testBlock.PrevBlockHash, prevHash, "prevhashes don't match")
	assert.Equal(t, testBlock.Jobs, jobs, "jobs don't match")
	assert.Equal(t, testBlock.MerkleHash, mHash, "merklehash doesn't match")
	assert.Nil(t, testBlock.Hash, "block hash is set")
}

func TestSetHash(t *testing.T) {
	jobs := []byte("test jobs")
	prevHash := []byte("00000000000000000000000000000000000000")
	mHash := []byte("0000000000000000000000000000000000000000")
	testBlock := NewBlock(jobs, prevHash, mHash)
	err := testBlock.SetHash()
	assert.NotNil(t, testBlock.Hash, "nil hash value")
	assert.Nil(t, err, "returned error")
	err = testBlock.SetHash()
	assert.NotNil(t, err, "didn't return error")
}

func TestVeriyBlock(t *testing.T) {
	jobs := []byte("test jobs")
	prevHash := []byte("00000000000000000000000000000000000000")
	mHash := []byte("0000000000000000000000000000000000000000")
	testBlock := NewBlock(jobs, prevHash, mHash)
	testBlock.SetHash()
	assert.True(t, testBlock.VerifyBlock(), "block failed verification")

	testBlock.Nonce = 50
	assert.False(t, testBlock.VerifyBlock(), "block passed verification")
}

func TestMarshalBlock(t *testing.T) {
	jobs := []byte("test jobs")
	prevHash := []byte("00000000000000000000000000000000000000")
	mHash := []byte("0000000000000000000000000000000000000(000")
	testBlock := NewBlock(jobs, prevHash, mHash)
	testBlock.SetHash()
	stringified, err := testBlock.MarshalBlock()
	var i interface{} = ""
	assert.Nil(t, err, "returned error")
	assert.IsType(t, i, stringified)

}
