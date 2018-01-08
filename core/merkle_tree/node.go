package merkle_tree

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"reflect"

	"github.com/gizo-network/gizo/helpers"
	"github.com/kpango/glg"
)

var Logger helpers.Log

type MerkleNode struct {
	Hash  []byte      `json:"hash"` //hash of a job struct
	Job   []byte      `json:"job"`
	Left  *MerkleNode `json:"left"`
	Right *MerkleNode `json:"right"`
}

// generates hash value of merklenode
func (n *MerkleNode) setHash() {
	l, err := n.Left.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	r, err := n.Right.Serialize()
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

//IsLeaf checks if the merklenode is a leaf node
func (n *MerkleNode) IsLeaf() bool {
	return n.Left.IsEmpty() && n.Right.IsEmpty()
}

//IsEmpty check if the merklenode is empty
func (n *MerkleNode) IsEmpty() bool {
	return reflect.ValueOf(n.Right).IsNil() && reflect.ValueOf(n.Left).IsNil() && reflect.ValueOf(n.Job).IsNil() && reflect.ValueOf(n.Hash).IsNil()
}

//IsEqual check if the input merklenode equals the merklenode calling the function
func (n MerkleNode) IsEqual(x MerkleNode) bool {
	nBytes, err := n.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	xBytes, err := x.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	return bytes.Equal(nBytes, xBytes)
}

//Serialize returns the bytes of a merklenode
func (x MerkleNode) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(x)
	return bytes, err
}

//NewNode returns a new merklenode
func NewNode(j []byte, lNode, rNode *MerkleNode) *MerkleNode {
	glg.Info("Creating MerkleNode")
	n := &MerkleNode{
		Left:  lNode,
		Right: rNode,
		Job:   j,
	}
	n.setHash()
	return n
}

//HashJobs hashes the jobs of two merklenodes
func HashJobs(x, y MerkleNode) []byte {
	headers := bytes.Join([][]byte{x.Job, y.Job}, []byte{})
	hash := sha256.Sum256(headers)
	return hash[:]
}
