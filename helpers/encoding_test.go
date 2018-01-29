package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/core/merkletree"
)

func TestEncode64(t *testing.T) {
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	b, err := node1.Serialize()
	assert.NoError(t, err)
	assert.NotNil(t, Encode64(b))
}

func TestDecode64(t *testing.T) {
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	b, err := node1.Serialize()
	assert.NoError(t, err)
	enc := Encode64(b)
	assert.Equal(t, b, Decode64(enc))
}
