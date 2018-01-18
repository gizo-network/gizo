package merkle_tree

import (
	"testing"

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
	node1 := NewNode([]byte("test1asdfasdf job"), &MerkleNode{}, &MerkleNode{})
	node2 := NewNode([]byte("test2 job asldkj;fasldkjfasd"), &MerkleNode{}, &MerkleNode{})
	node3 := NewNode([]byte("test3 asdfasl;dfasdjob"), &MerkleNode{}, &MerkleNode{})
	node4 := NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node5 := NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node6 := NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node7 := NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &MerkleNode{}, &MerkleNode{})
	node8 := NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}

	tree := MerkleTree{
		LeafNodes: nodes,
	}
	assert.Nil(t, tree.GetRoot())
	tree.Build()
	assert.NotNil(t, tree.GetRoot())
}

func TestNewMerkleTree(t *testing.T) {
	node1 := NewNode([]byte("test1asdfasdf job"), &MerkleNode{}, &MerkleNode{})
	node2 := NewNode([]byte("test2 job asldkj;fasldkjfasd"), &MerkleNode{}, &MerkleNode{})
	node3 := NewNode([]byte("test3 asdfasl;dfasdjob"), &MerkleNode{}, &MerkleNode{})
	node4 := NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node5 := NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node6 := NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node7 := NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &MerkleNode{}, &MerkleNode{})
	node8 := NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}

	tree := NewMerkleTree(nodes)
	assert.NotNil(t, tree.GetRoot())
	assert.NotNil(t, tree.GetLeafNodes())
}

func TestVerifyTree(t *testing.T) {
	node1 := NewNode([]byte("test1asdfasdf job"), &MerkleNode{}, &MerkleNode{})
	node2 := NewNode([]byte("test2 job asldkj;fasldkjfasd"), &MerkleNode{}, &MerkleNode{})
	node3 := NewNode([]byte("test3 asdfasl;dfasdjob"), &MerkleNode{}, &MerkleNode{})
	node4 := NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node5 := NewNode([]byte("tesasdfa;sdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node6 := NewNode([]byte("tesasdfa;sadasdfasdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	node7 := NewNode([]byte("tesasdfa;sdlkfj;asasdfasfdat4 job"), &MerkleNode{}, &MerkleNode{})
	node8 := NewNode([]byte("tesasdfasdfsadfasdfa;sdlkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}

	tree := NewMerkleTree(nodes)
	assert.True(t, tree.VerifyTree())

	tree.SetLeafNodes(tree.GetLeafNodes()[2:])
	assert.False(t, tree.VerifyTree())
}
