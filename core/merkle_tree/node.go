package merkle_tree

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"log"
)

type Node struct {
	Hash   []byte `json:"hash"` //hash of a job struct
	Job    []byte `json:"job"`
	Parent *Node  `json:"parent"`
	Left   *Node  `json:"left"`
	Right  *Node  `json:"right"`
}

func NewNode(j []byte, pNode, lNode, rNode *Node) *Node {
	return &Node{
		Parent: pNode,
		Left:   lNode,
		Right:  rNode,
		Job:    j,
	}
}

func HashJobs(x, y Node) []byte {
	headers := bytes.Join([][]byte{x.Job, y.Job}, []byte{})
	hash := sha256.Sum256(headers)
	return hash[:]
}

func MarshalNode(x Node) ([]byte, error) {
	bytes, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func EmptyNode(n Node) bool {
	return n.Right == nil && n.Left == nil && n.Parent == nil && n.Job == nil
}

func (n *Node) SetHash() {
	p, err := MarshalNode(*n.Parent)
	if err != nil {
		log.Fatal(err)
	}
	l, err := MarshalNode(*n.Left)
	if err != nil {
		log.Fatal(err)
	}
	r, err := MarshalNode(*n.Right)
	if err != nil {
		log.Fatal(err)
	}

	headers := bytes.Join([][]byte{p, l, r, n.Job}, []byte{})
	if err != nil {
		log.Fatal(err)
	}
	hash := sha256.Sum256(headers)
	n.Hash = hash[:]
}

func (n *Node) SetParent(p *Node) {
	n.Parent = p
}

func (n *Node) IsLeaf() bool {
	return EmptyNode(*n.Left) && EmptyNode(*n.Right)
}

func (n *Node) IsRoot() bool {
	return EmptyNode(*n.Parent)
}

func (n *Node) IsEmpty() bool {
	return EmptyNode(*n.Parent) && EmptyNode(*n.Left) && EmptyNode(*n.Right)
}
