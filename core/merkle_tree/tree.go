package merkle_tree

import (
	"bytes"
	"errors"
	"sync"

	"github.com/gizo-network/gizo/core"
	"github.com/kpango/glg"
)

type MerkleTree struct {
	Root      *MerkleNode   `json:"root"`
	LeafNodes []*MerkleNode `json:"leafNodes"`
}

var ErrTooMuchLeafNodes = errors.New("core/merkle tree: length of leaf nodes is greater than 24")
var ErrOddLeafNodes = errors.New("core/merkle tree: odd number of leaf nodes")
var ErrTreeRebuildAttempt = errors.New("core/merkle tree: attempt to rebuild tree")
var ErrTreeNotBuilt = errors.New("core/merkle_tree: tree hasn't been built")

// NewMerkleTree returns empty merkletree
func NewMerkleTree(nodes []*MerkleNode) *MerkleTree {
	return &MerkleTree{
		LeafNodes: nodes,
	}
}

func merge(left, right MerkleNode) *MerkleNode {
	parent := NewNode(HashJobs(left, right), &left, &right)
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
				shrink = append(shrink, shrink[len(shrink)-1]) //duplicate last  to balance tree
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

func (m *MerkleTree) BreakTree() []MerkleNode {
	var mutex sync.Mutex
	leafs := []MerkleNode{}
	queue := make(chan MerkleNode, 100)
	// done := make(chan bool)

	if m.Root.IsEmpty() {
		glg.Fatal(ErrTreeNotBuilt)
	} else {
		glg.Info("deploying root children")

		queue <- *m.Root.Left
		queue <- *m.Root.Right

		for len(queue) != 0 {
			select {
			case node := <-queue:
				// nBytes, _ := MarshalNode(node)
				// b, _ := helpers.PrettyJSON(nBytes)
				glg.Info("received node")
				if node.IsLeaf() {
					glg.Info("detected leaf")
					mutex.Lock()
					leafs = append(leafs, node)
					mutex.Unlock()
				} else {
					queue <- *node.Left
					queue <- *node.Right
				}
			default:
				// close(queue)
				// fmt.Println("empty queue")
				break
			}
		}
	}
	// queue <- *m.Root.Left
	// queue <- *m.Root.Right

	// go func() {
	// 	time.Sleep(time.Second * 2)
	// 	queue <- *m.Root.Left
	// 	queue <- *m.Root.Right
	// }()
	// for {
	// 	select {
	// 	case node := <-queue:
	// 		fmt.Println(node)
	// 	}
	// }

	return leafs
}

//VerifyTree returns true if tree is verified
func (m MerkleTree) VerifyTree() bool {
	t := NewMerkleTree(m.LeafNodes)
	err := t.BuildTree()
	if err != nil {
		glg.Fatal(err)
	}
	mBytes, err := MarshalNode(*m.Root)
	if err != nil {
		glg.Fatal(err)
	}
	tBytes, err := MarshalNode(*t.Root)
	if err != nil {
		glg.Fatal(err)
	}
	return bytes.Compare(tBytes, mBytes) == 0
}

//SearchTree returns true if node with has exists
// func SearchTree(hash []byte) bool {

// }
