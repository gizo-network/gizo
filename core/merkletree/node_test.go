package merkletree_test

import (
	"encoding/hex"
	"testing"

	"github.com/gizo-network/gizo/crypt"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	n := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	assert.NotNil(t, n.GetHash(), "empty hash value")
	assert.NotNil(t, n, "returned empty node")
}

func TestMarshalMerkleNode(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	n := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	b, err := n.Serialize()
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestIsLeaf(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	n := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	assert.True(t, n.IsLeaf())
}

func TestIsEmpty(t *testing.T) {
	n := merkletree.MerkleNode{}
	assert.True(t, n.IsEmpty())
}

func TestIsEqual(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	n := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	assert.True(t, n.IsEqual(*n))
}
