package merkle_tree

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"reflect"

	"github.com/kpango/glg"
)

type MerkleNode struct {
	Hash  []byte      `json:"hash"` //hash of a job struct
	Job   []byte      `json:"job"`
	Left  *MerkleNode `json:"left"`
	Right *MerkleNode `json:"right"`
}

func NewNode(j []byte, lNode, rNode *MerkleNode) *MerkleNode {
	return &MerkleNode{
		Left:  lNode,
		Right: rNode,
		Job:   j,
	}
}

func HashJobs(x, y MerkleNode) []byte {
	headers := bytes.Join([][]byte{x.Job, y.Job}, []byte{})
	hash := sha256.Sum256(headers)
	return hash[:]
}

func MarshalNode(x MerkleNode) ([]byte, error) {
	bytes, err := json.Marshal(x)
	return bytes, err
}

func (n *MerkleNode) SetHash() {
	l, err := MarshalNode(*n.Left)
	if err != nil {
		glg.Fatal(err)
	}
	r, err := MarshalNode(*n.Right)
	if err != nil {
		glg.Fatal(err)
	}

	headers := bytes.Join([][]byte{l, r, n.Job}, []byte{})
	if err != nil {
		glg.Fatal(err)
	}
	hash := sha256.Sum256(headers)
	n.Hash = hash[:]
}

// func (n *MerkleNode) SetParent(p *MerkleNode) {
// 	n.Parent = *p
// }

func (n *MerkleNode) IsLeaf() bool {
	return n.Left.IsEmpty() && n.Right.IsEmpty()
}

func (n *MerkleNode) IsRoot() bool {
	return n.IsEmpty() == false && n.IsLeaf() == false
}

func (n *MerkleNode) IsEmpty() bool {
	return reflect.ValueOf(n.Right).IsNil() && reflect.ValueOf(n.Left).IsNil() && reflect.ValueOf(n.Job).IsNil() && reflect.ValueOf(n.Hash).IsNil()
}
