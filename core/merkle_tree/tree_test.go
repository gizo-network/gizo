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
	assert.NotNil(t, ErrLeafNodesNotEmpty)
}

func BenchmarkDismantle4Nodes(b *testing.B) {
	node1 := NewNode([]byte("test1asdfasdf job"), &MerkleNode{}, &MerkleNode{})
	node2 := NewNode([]byte("test2 job asldkj;fasldkjfasd"), &MerkleNode{}, &MerkleNode{})
	node3 := NewNode([]byte("test3 asdfasl;dfasdjob"), &MerkleNode{}, &MerkleNode{})
	node4 := NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &MerkleNode{}, &MerkleNode{})
	nodes := []*MerkleNode{node1, node2, node3, node4}

	tree := NewMerkleTree(nodes)

	for i := 0; i < b.N; i++ {
		newTree := MerkleTree{}
		newTree.Root = tree.Root
		newTree.Dismantle()
	}
}

func BenchmarkDismantle8Nodes(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		newTree := MerkleTree{}
		newTree.Root = tree.Root
		newTree.Dismantle()
	}
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
	assert.Nil(t, tree.Root)
	tree.Build()
	assert.NotNil(t, tree.Root)
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
	assert.NotNil(t, tree.Root)
	assert.NotNil(t, tree.LeafNodes)
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

	tree.LeafNodes[0] = &MerkleNode{}
	assert.False(t, tree.VerifyTree())
}

func TestDismantle(t *testing.T) {
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

	newTree := &MerkleTree{
		Root: tree.Root,
	}
	newTree.Dismantle()
	assert.NotNil(t, newTree.LeafNodes)
	for _, val := range tree.LeafNodes {
		assert.True(t, newTree.SearchLeaf(val.Hash))
	}
}
