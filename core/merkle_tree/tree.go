package merkle_tree

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/kpango/glg"
)

var ErrTooMuchLeafNodes = errors.New("core/merkle tree: length of leaf nodes is greater than 24")
var ErrOddLeafNodes = errors.New("core/merkle tree: odd number of leaf nodes")
var ErrTreeRebuildAttempt = errors.New("core/merkle tree: attempt to rebuild tree")
var ErrTreeNotBuilt = errors.New("core/merkle_tree: tree hasn't been built")
var ErrLeafNodesEmpty = errors.New("core/merkle_tree: leafnodes is empty")

type MerkleTree struct {
	Root      []byte        `json:"root"`
	LeafNodes []*MerkleNode `json:"leafNodes"`
}

//builds merkle tree from leafs to root and sets the root of the merkletree
func (m *MerkleTree) Build() error {
	glg.Info("Building MerkleTree")
	if reflect.ValueOf(m.Root).IsNil() == false {
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
		m.Root = shrink[0].Hash
	}
	return nil
}

//Serialize returns the bytes of a merkletree
func (x MerkleTree) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(x)
	return bytes, err
}

//VerifyTree returns true if tree is verified
func (m MerkleTree) VerifyTree() bool {
	t := NewMerkleTree(m.LeafNodes)
	return bytes.Equal(t.Root, m.Root)
}

// Search returns true if node with hash exists
func (m MerkleTree) Search(hash []byte) (bool, error) {
	if len(m.LeafNodes) == 0 {
		return false, ErrLeafNodesEmpty
	} else {
		for _, n := range m.LeafNodes {
			if bytes.Equal(n.Hash, hash) {
				return true, nil
			}
		}
	}
	return false, nil
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
