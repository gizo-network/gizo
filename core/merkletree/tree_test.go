package merkletree

import (
	"testing"

	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert.NotNil(t, ErrTooMuchLeafNodes)
	assert.NotNil(t, ErrOddLeafNodes)
	assert.NotNil(t, ErrTreeRebuildAttempt)
	assert.NotNil(t, ErrTreeNotBuilt)
	assert.NotNil(t, ErrLeafNodesEmpty)
}

func TestBuild(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}; test()")
	node1 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node2 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node3 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node4 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node5 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node6 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node7 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node8 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}

	tree := MerkleTree{
		LeafNodes: nodes,
	}
	assert.Nil(t, tree.GetRoot())
	tree.Build()
	assert.NotNil(t, tree.GetRoot())
}

func TestNewMerkleTree(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}; test()")
	node1 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node2 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node3 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node4 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node5 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node6 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node7 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node8 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}

	tree := NewMerkleTree(nodes)
	assert.NotNil(t, tree.GetRoot())
	assert.NotNil(t, tree.GetLeafNodes())
}

func TestVerifyTree(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}; test()")
	node1 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node2 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node3 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node4 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node5 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node6 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node7 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node8 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}

	tree := NewMerkleTree(nodes)
	assert.True(t, tree.VerifyTree())

	tree.SetLeafNodes(tree.GetLeafNodes()[4:])
	assert.False(t, tree.VerifyTree())
}

func TestSearchNode(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}; test()")
	node1 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node2 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node3 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node4 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node5 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node6 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node7 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node8 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}

	tree := NewMerkleTree(nodes)
	f, err := tree.SearchNode(node5.GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestSearchJob(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}; test()")
	node1 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node2 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node3 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node4 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node5 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node6 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node7 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node8 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}

	tree := NewMerkleTree(nodes)
	f, err := tree.SearchJob(j.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestMerge(t *testing.T) {
	j := job.NewJob("func test(){return 1+1}; test()")
	node1 := NewNode(*j, &MerkleNode{}, &MerkleNode{})
	node2 := NewNode(*j, &MerkleNode{}, &MerkleNode{})

	parent := merge(*node1, *node2)
	assert.NotNil(t, parent)
	assert.NotNil(t, parent.GetHash())
	assert.Equal(t, node1.GetHash(), parent.GetLeftNode().GetHash())
	assert.Equal(t, node2.GetHash(), parent.GetRightNode().GetHash())
}
