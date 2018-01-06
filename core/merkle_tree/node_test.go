package merkle_tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	n := NewNode([]byte("test job"), &MerkleNode{}, &MerkleNode{})
	assert.NotNil(t, n.Hash, "empty hash value")
	assert.NotNil(t, n, "returned empty node")
}

func TestHashJobs(t *testing.T) {
	l := NewNode([]byte("test job 1"), &MerkleNode{}, &MerkleNode{})
	r := NewNode([]byte("test job 2"), &MerkleNode{}, &MerkleNode{})
	b := HashJobs(*l, *r)
	assert.NotNil(t, b)
}

func TestMarshalMerkleNode(t *testing.T) {
	n := NewNode([]byte("test job"), &MerkleNode{}, &MerkleNode{})
	b, err := n.Serialize()
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestIsLeaf(t *testing.T) {
	n := NewNode([]byte("test job"), &MerkleNode{}, &MerkleNode{})
	assert.True(t, n.IsLeaf())
}

func TestIsEmpty(t *testing.T) {
	n := MerkleNode{}
	assert.True(t, n.IsEmpty())
}

func TestIsEqual(t *testing.T) {
	n := NewNode([]byte("test job"), &MerkleNode{}, &MerkleNode{})
	assert.True(t, n.IsEqual(*n))
}
