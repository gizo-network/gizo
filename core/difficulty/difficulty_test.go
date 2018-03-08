package difficulty_test

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/gizo-network/gizo/crypt"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/difficulty"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
)

func TestDifficulty(t *testing.T) {
	os.Setenv("ENV", "dev")
	core.RemoveDataPath()
	bc := core.CreateBlockChain()
	priv, _ := crypt.GenKeys()
	j := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	node1 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})

	nodes := []*merkletree.MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8}
	tree := merkletree.NewMerkleTree(nodes)
	block := core.NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestHeight(), 10)
	bc.AddBlock(block)
	d10 := benchmark.NewBenchmark(0.0115764096, 10)
	d11 := benchmark.NewBenchmark(0.13054728, 11)
	d12 := benchmark.NewBenchmark(0.0740971, 12)
	d13 := benchmark.NewBenchmark(0.28987127999999995, 13)
	d14 := benchmark.NewBenchmark(1.36593388, 14)
	d15 := benchmark.NewBenchmark(1.8645611, 15)
	d16 := benchmark.NewBenchmark(3.82076494, 16)
	d17 := benchmark.NewBenchmark(7.12966816, 17)
	d18 := benchmark.NewBenchmark(28.470944839999998, 18)
	d19 := benchmark.NewBenchmark(42.251310620000005, 19)
	benchmarks := []benchmark.Benchmark{d10, d11, d12, d13, d14, d15, d16, d17, d18, d19}
	assert.NotNil(t, difficulty.Difficulty(benchmarks, *bc))
}
