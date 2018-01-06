package core

import (
	"testing"

	"github.com/gizo-network/gizo/core/merkle_tree"
	"github.com/stretchr/testify/assert"
)

type String string

func TestNewBlock(t *testing.T) {
	node1 := merkle_tree.NewNode([]byte("test1asdfasdf job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node2 := merkle_tree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node3 := merkle_tree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node4 := merkle_tree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	nodes := []*merkle_tree.MerkleNode{node1, node2, node3, node4}
	tree := merkle_tree.NewMerkleTree(nodes)
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0)

	assert.NotNil(t, testBlock, "returned empty tblock")
	assert.Equal(t, testBlock.PrevBlockHash, prevHash, "prevhashes don't match")
}

func TestVeriyBlock(t *testing.T) {
	node1 := merkle_tree.NewNode([]byte("test1asdfasdf job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node2 := merkle_tree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node3 := merkle_tree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node4 := merkle_tree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	nodes := []*merkle_tree.MerkleNode{node1, node2, node3, node4}
	tree := merkle_tree.NewMerkleTree(nodes)
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0)

	assert.True(t, testBlock.VerifyBlock(), "block failed verification")

	testBlock.Nonce = 50
	assert.False(t, testBlock.VerifyBlock(), "block passed verification")
}

func TestSerialize(t *testing.T) {
	node1 := merkle_tree.NewNode([]byte("test1asdfasdf job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node2 := merkle_tree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node3 := merkle_tree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node4 := merkle_tree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	nodes := []*merkle_tree.MerkleNode{node1, node2, node3, node4}
	tree := merkle_tree.NewMerkleTree(nodes)
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0)
	stringified, err := testBlock.Serialize()
	assert.Nil(t, err, "returned error")
	assert.NotEmpty(t, stringified)
}

func TestUnMarshalBlock(t *testing.T) {
	node1 := merkle_tree.NewNode([]byte("test1asdfasdf job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node2 := merkle_tree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node3 := merkle_tree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	node4 := merkle_tree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkle_tree.MerkleNode{}, &merkle_tree.MerkleNode{})
	nodes := []*merkle_tree.MerkleNode{node1, node2, node3, node4}
	tree := merkle_tree.NewMerkleTree(nodes)
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0)
	stringified, _ := testBlock.Serialize()
	unmarshaled, err := DeserilizeBlock(stringified)
	assert.Nil(t, err)
	assert.Equal(t, testBlock, unmarshaled)
}
