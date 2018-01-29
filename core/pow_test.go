package core

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/core/merkletree"
)

func TestNewPOW(t *testing.T) {
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	nodes := []*merkletree.MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}
	tree := merkletree.NewMerkleTree(nodes)

	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: []byte("00000000000000000000000000000000000000"),
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(10)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: 0,
	}
	pow := NewPOW(block)
	assert.NotNil(t, pow)
}

func TestRun(t *testing.T) {
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	nodes := []*merkletree.MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}
	tree := merkletree.NewMerkleTree(nodes)

	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: []byte("00000000000000000000000000000000000000"),
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(10)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: 0,
	}
	pow := NewPOW(block)
	assert.NotNil(t, pow)
	pow.Run()
	assert.NotNil(t, block.GetHeader().GetHash())
	assert.NotNil(t, block.GetHeader().GetNonce())
}

func TestPrepareData(t *testing.T) {
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	nodes := []*merkletree.MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}
	tree := merkletree.NewMerkleTree(nodes)

	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: []byte("00000000000000000000000000000000000000"),
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(10)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: 0,
	}
	pow := NewPOW(block)
	assert.NotNil(t, pow)
	assert.NotNil(t, pow.prepareData(5))
}

func TestValidate(t *testing.T) {
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	nodes := []*merkletree.MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}
	tree := merkletree.NewMerkleTree(nodes)

	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: []byte("00000000000000000000000000000000000000"),
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(10)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: 1,
	}
	pow := NewPOW(block)
	pow.Run()
	assert.NotNil(t, pow)
	assert.True(t, pow.Validate())
}
