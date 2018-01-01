package merkle_tree

import (
	"errors"
	"log"

	"github.com/gizo-network/gizo/core"
)

type MerkleTree struct {
	Root      *Node   `json:"root"`
	LeafNodes []*Node `json:"leafNodes"`
}

var ErrTooMuchLeafNodes = errors.New("merkle tree: length of leaf nodes is greater than 24")
var ErrOddLeafNodes = errors.New("merkle tree: odd number of leaf nodes")
var ErrTreeRebuildAttempt = errors.New("merkle tree: attempt to rebuild tree")

// NewMerkleTree returns empty merkletree
func NewMerkleTree(nodes []*Node) *MerkleTree {
	return &MerkleTree{
		LeafNodes: nodes,
	}
}

func Merge(left, right *Node) *Node {
	parent := NewNode(HashJobs(*left, *right), &Node{}, left, right)
	parent.SetHash()
	return parent
}

func (m *MerkleTree) BuildTree() error {
	//FIXME: add parent to nodes
	if m.Root != nil {
		return ErrTreeRebuildAttempt
	}
	if int64(len(m.LeafNodes)) > core.MaxTreeJobs.Int64() {
		return ErrTooMuchLeafNodes
	} else if len(m.LeafNodes)%2 != 0 {
		return ErrOddLeafNodes
	} else {
		var shrink []*Node = m.LeafNodes
		for len(shrink) != 1 {
			var temp []*Node
			if len(shrink)%2 == 0 {
				for i := 0; i < len(shrink); i += 2 {
					parent := Merge(shrink[i], shrink[i+1])
					temp = append(temp, parent)
				}
			} else {
				shrink = append(shrink, shrink[len(shrink)-1]) //duplicate last  to balance tree
				for i := 0; i < len(shrink); i += 2 {
					parent := Merge(shrink[i], shrink[i+1])
					temp = append(temp, parent)
				}
			}
			shrink = temp
		}
		m.Root = shrink[0]
	}
	return nil
}

//VerifyTree returns true if tree is verified
func (m *MerkleTree) VerifyTree() bool {
	t := NewMerkleTree(m.LeafNodes)
	err := t.BuildTree()
	if err != nil {
		log.Fatal(err)
	}

	mBytes, err := MarshalNode(*m.Root)
	if err != nil {
		log.Fatal(err)
	}
	tBytes, err := MarshalNode(*t.Root)
	if err != nil {
		log.Fatal(err)
	}
	return string(tBytes) == string(mBytes)
}

//SearchHash returns true if node with has exists
func SearchHash(hash []byte) bool {

}
