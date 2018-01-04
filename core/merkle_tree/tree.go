package merkle_tree

import (
	"bytes"
	"errors"
	"sync"

	"github.com/kpango/glg"
)

var ErrTooMuchLeafNodes = errors.New("core/merkle tree: length of leaf nodes is greater than 24")
var ErrOddLeafNodes = errors.New("core/merkle tree: odd number of leaf nodes")
var ErrTreeRebuildAttempt = errors.New("core/merkle tree: attempt to rebuild tree")
var ErrTreeNotBuilt = errors.New("core/merkle_tree: tree hasn't been built")
var ErrLeafNodesNotEmpty = errors.New("core/merkle_tree: leafnodes is not empty")

type MerkleTree struct {
	Root      *MerkleNode   `json:"root"`
	LeafNodes []*MerkleNode `json:"leafNodes"`
}

//builds merkle tree from leafs to root and sets the root of the merkletree
func (m *MerkleTree) Build() error {
	glg.Info("Building MerkleTree")
	if m.Root != nil {
		return ErrTreeRebuildAttempt
	}
	if int64(len(m.LeafNodes)) > MaxTreeJobs.Int64() {
		return ErrTooMuchLeafNodes
	} else if len(m.LeafNodes)%2 != 0 {
		return ErrOddLeafNodes
	} else {
		var shrink []*MerkleNode = m.LeafNodes
		for len(shrink) != 1 {
			var levelUp []*MerkleNode
			if len(shrink)%2 == 0 {
				for i := 0; i < len(shrink); i += 2 {
					parent := merge(*shrink[i], *shrink[i+1])
					levelUp = append(levelUp, parent)
				}
			} else {
				glg.Warn("core/merkle_tree: Duplicating solo node...")
				shrink = append(shrink, shrink[len(shrink)-1]) //duplicate last to balance tree
				for i := 0; i < len(shrink); i += 2 {
					parent := merge(*shrink[i], *shrink[i+1])
					levelUp = append(levelUp, parent)
				}
			}
			shrink = levelUp
		}
		m.Root = shrink[0]
	}
	return nil
}

//Dismantle breaks down a root into it's leaves
func (m *MerkleTree) Dismantle() {
	//FIXME: return leafs in original sequence
	glg.Info("Dismantling MerkleNode")
	var mutex sync.Mutex
	var wg sync.WaitGroup
	leafs := []*MerkleNode{}
	queue := make(chan MerkleNode, 100)

	if m.Root.IsEmpty() {
		glg.Fatal(ErrTreeNotBuilt)
	} else if len(m.LeafNodes) != 0 {
		glg.Fatal(ErrLeafNodesNotEmpty)
	} else {
		queue <- *m.Root
		for len(queue) != 0 {
			select {
			case node := <-queue:
				if node.IsLeaf() && node.IsEmpty() == false {
					wg.Add(1)
					go func() {
						mutex.Lock()
						leafs = append(leafs, &node)
						mutex.Unlock()
						wg.Done()
					}()
				} else {
					queue <- *node.Left
					queue <- *node.Right
				}
			}
		}
		wg.Wait()
	}
	m.LeafNodes = stripDuplicates(leafs)
}

//VerifyTree returns true if tree is verified
func (m MerkleTree) VerifyTree() bool {
	t := NewMerkleTree(m.LeafNodes)
	mBytes, err := MarshalMerkleNode(*m.Root)
	if err != nil {
		glg.Fatal(err)
	}
	tBytes, err := MarshalMerkleNode(*t.Root)
	if err != nil {
		glg.Fatal(err)
	}
	return bytes.Equal(tBytes, mBytes)
}

//SearchTree returns true if node with has exists
func (m MerkleTree) SearchLeaf(hash []byte) bool {
	if len(m.LeafNodes) == 0 {
		m.Dismantle()
		for _, n := range m.LeafNodes {
			if bytes.Equal(n.Hash, hash) {
				return true
			}
		}
	} else {
		for _, n := range m.LeafNodes {
			if bytes.Equal(n.Hash, hash) {
				return true
			}
		}
	}
	return false
}

// NewMerkleTree returns empty merkletree
func NewMerkleTree(nodes []*MerkleNode) *MerkleTree {
	glg.Info("Creating MerkleTree")
	t := &MerkleTree{
		LeafNodes: nodes,
	}
	err := t.Build()
	if err != nil {
		glg.Fatal(err)
	}
	return t
}

//merges two nodes
func merge(left, right MerkleNode) *MerkleNode {
	parent := NewNode(HashJobs(left, right), &left, &right)
	return parent
}

//removes the duplicates from an array of merklenodes
func stripDuplicates(input []*MerkleNode) []*MerkleNode {
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input); j++ {
			if i == j {
				continue
			}
			if input[i].IsEqual(*input[j]) {
				input = append(input[:j], input[j+1:]...)
			}
		}
	}
	return input
}
