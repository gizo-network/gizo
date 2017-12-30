package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlock(t *testing.T) {
	jobs := []byte("test jobs")
	prevHash := []byte("00000000000000000000000000000000000000")
	mHash := []byte("0000000000000000000000000000000000000000")
	testBlock := NewBlock(jobs, prevHash, mHash)

	assert.NotNil(t, testBlock)
	assert.Equal(t, testBlock.PrevBlockHash, prevHash)
	assert.Equal(t, testBlock.Jobs, jobs)
	assert.Equal(t, testBlock.MerkleHash, mHash)
	assert.Nil(t, testBlock.Hash)
}

func TestSetHash(t *testing.T) {
	jobs := []byte("test jobs")
	prevHash := []byte("00000000000000000000000000000000000000")
	mHash := []byte("0000000000000000000000000000000000000000")
	testBlock := NewBlock(jobs, prevHash, mHash)
	err := testBlock.SetHash()
	assert.NotNil(t, testBlock.Hash)
	assert.Nil(t, err)
	err = testBlock.SetHash()
	assert.NotNil(t, err)
}

func TestVeriyBlock(t *testing.T) {
	jobs := []byte("test jobs")
	prevHash := []byte("00000000000000000000000000000000000000")
	mHash := []byte("0000000000000000000000000000000000000000")
	testBlock := NewBlock(jobs, prevHash, mHash)
	testBlock.SetHash()
	assert.True(t, testBlock.VerifyBlock())

	testBlock.Nonce = 50
	assert.False(t, testBlock.VerifyBlock())
}
