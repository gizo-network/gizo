package merkletree

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/kpango/glg"
)

var (
	ErrNodeDoesntExist    = errors.New("core/merkletree: node doesn't exist")
	ErrLeafNodesEmpty     = errors.New("core/merkletree: leafnodes is empty")
	ErrTreeNotBuilt       = errors.New("core/merkletree: tree hasn't been built")
	ErrTreeRebuildAttempt = errors.New("core/merkle tree: attempt to rebuild tree")
	ErrOddLeafNodes       = errors.New("core/merkle tree: odd number of leaf nodes")
	ErrTooMuchLeafNodes   = errors.New("core/merkle tree: length of leaf nodes is greater than 24")
)

// MerkleTree tree of jobs
type MerkleTree struct {
	Root      []byte        `json:"root"`
	LeafNodes []*MerkleNode `json:"leafNodes"`
}

// GetRoot returns root
func (m MerkleTree) GetRoot() []byte {
	return m.Root
}

func (m *MerkleTree) setRoot(r []byte) {
	m.Root = r
}

// GetLeafNodes return leafnodes
func (m MerkleTree) GetLeafNodes() []*MerkleNode {
	return m.LeafNodes
}

// SetLeafNodes return leafnodes
func (m *MerkleTree) SetLeafNodes(l []*MerkleNode) {
	m.LeafNodes = l
}

//Build builds merkle tree from leafs to root, hashed the root and sets it as the root of the merkletree
func (m *MerkleTree) Build() error {
	if reflect.ValueOf(m.GetRoot()).IsNil() == false {
		return ErrTreeRebuildAttempt
	}
	if int64(len(m.GetLeafNodes())) > MaxTreeJobs.Int64() {
		return ErrTooMuchLeafNodes
	} else if len(m.GetLeafNodes())%2 != 0 {
		return ErrOddLeafNodes
	} else {
		var shrink = m.LeafNodes
		for len(shrink) != 1 {
			var levelUp []*MerkleNode
			if len(shrink)%2 == 0 {
				for i := 0; i < len(shrink); i += 2 {
					parent := merge(*shrink[i], *shrink[i+1])
					levelUp = append(levelUp, parent)
				}
			} else {
				glg.Warn("core/merkletree: Duplicating solo node...")
				shrink = append(shrink, shrink[len(shrink)-1]) //duplicate last to balance tree
				for i := 0; i < len(shrink); i += 2 {
					parent := merge(*shrink[i], *shrink[i+1])
					levelUp = append(levelUp, parent)
				}
			}
			shrink = levelUp
		}
		m.Root = shrink[0].GetHash()
	}
	return nil
}

//Serialize returns the bytes of a merkletree
func (m MerkleTree) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	return bytes, err
}

//VerifyTree returns true if tree is verified
func (m MerkleTree) VerifyTree() bool {
	t := NewMerkleTree(m.GetLeafNodes())
	return bytes.Equal(t.GetRoot(), m.GetRoot())
}

// Search returns true if node with hash exists
func (m MerkleTree) Search(hash []byte) (*MerkleNode, error) {
	if len(m.GetLeafNodes()) == 0 {
		return nil, ErrLeafNodesEmpty
	}
	for _, n := range m.GetLeafNodes() {
		if bytes.Equal(n.GetHash(), hash) {
			return n, nil
		}
	}
	return nil, ErrNodeDoesntExist
}

// NewMerkleTree returns empty merkletree
func NewMerkleTree(nodes []*MerkleNode) *MerkleTree {
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
