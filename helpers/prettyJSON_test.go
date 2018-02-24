package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job"
)

func TestPrettyJSON(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}", "test")
	node1 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	b, err := node1.Serialize()
	assert.NoError(t, err)
	pretty, err := helpers.PrettyJSON(b)
	assert.NoError(t, err)
	assert.NotNil(t, pretty)
}
