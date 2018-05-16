package merkletree

import (
	"encoding/hex"
	"testing"

	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	node1 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node2 := NewNode(*j, &MerkleNode{}, &MerkleNode{})

	parent := merge(*node1, *node2)
	assert.NotNil(t, parent)
	assert.NotNil(t, parent.GetHash())
	assert.Equal(t, node1.GetHash(), parent.GetLeftNode().GetHash())
	assert.Equal(t, node2.GetHash(), parent.GetRightNode().GetHash())
}
