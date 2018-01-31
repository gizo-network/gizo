package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/core/merkletree"
)

func TestPrettyJSON(t *testing.T) {
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	b, err := node1.Serialize()
	assert.NoError(t, err)
	pretty, err := PrettyJSON(b)
	assert.NoError(t, err)
	assert.NotNil(t, pretty)
}
