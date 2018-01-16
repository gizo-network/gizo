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
	hash  []byte      `json:"hash"` //hash of a job struct
	job   []byte      `json:"job"`
	left  *MerkleNode `json:"left"`
	right *MerkleNode `json:"right"`
}

func (n MerkleNode) GetHash() []byte {
	return n.hash
}

// generates hash value of merklenode
func (n *MerkleNode) setHash() {
	l, err := n.left.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	r, err := n.right.Serialize()
	if err != nil {
		glg.Fatal(err)
	}

	headers := bytes.Join([][]byte{l, r, n.job}, []byte{})
	if err != nil {
		glg.Fatal(err)
	}
	hash := sha256.Sum256(headers)
	n.hash = hash[:]
}

func (n MerkleNode) GetJob() []byte {
	return n.job
}

func (n *MerkleNode) SetJob(j []byte) {
	n.job = j
}

func (n MerkleNode) GetLeftNode() MerkleNode {
	return *n.left
}

func (n *MerkleNode) SetLeftNode(l MerkleNode) {
	n.left = &l
}

func (n MerkleNode) GetRightNOde() MerkleNode {
	return *n.right
}

func (n *MerkleNode) SetRightNode(r MerkleNode) {
	n.right = &r
}

//IsLeaf checks if the merklenode is a leaf node
func (n *MerkleNode) IsLeaf() bool {
	return n.left.IsEmpty() && n.right.IsEmpty()
}

//IsEmpty check if the merklenode is empty
func (n *MerkleNode) IsEmpty() bool {
	return reflect.ValueOf(n.right).IsNil() && reflect.ValueOf(n.left).IsNil() && reflect.ValueOf(n.job).IsNil() && reflect.ValueOf(n.hash).IsNil()
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
		left:  lNode,
		right: rNode,
		job:   j,
	}
	n.setHash()
	return n
}

//HashJobs hashes the jobs of two merklenodes
func HashJobs(x, y MerkleNode) []byte {
	headers := bytes.Join([][]byte{x.job, y.job}, []byte{})
	hash := sha256.Sum256(headers)
	return hash[:]
}
