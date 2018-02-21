package merkletree

import (
	"testing"

	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}", "test")
	n := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	assert.NotNil(t, n.GetHash(), "empty hash value")
	assert.NotNil(t, n, "returned empty node")
}

func TestMarshalMerkleNode(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}", "test")
	n := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	b, err := n.Serialize()
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestIsLeaf(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}", "test")
	n := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	assert.True(t, n.IsLeaf())
}

func TestIsEmpty(t *testing.T) {
	n := MerkleNode{}
	assert.True(t, n.IsEmpty())
}

func TestIsEqual(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}", "test")
	n := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	assert.True(t, n.IsEqual(*n))
}
