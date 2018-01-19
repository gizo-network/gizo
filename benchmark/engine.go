package benchmark

import (
	"math/rand"
	"time"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
)

type BenchmarkEngine struct {
	Data  []Benchmark
	Score uint8
}

func (b *BenchmarkEngine) SetScore(s uint8) {
	b.Score = s
}

func (b BenchmarkEngine) GetScore() uint8 {
	return b.Score
}

func (b *BenchmarkEngine) AddBenchmark(benchmark Benchmark) {
	b.Data = append(b.Data, benchmark)
}

func (b BenchmarkEngine) GetData() []Benchmark {
	return b.Data
}

func (b BenchmarkEngine) Block(difficulty uint8) *core.Block {
	//random data
	node1 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node9 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node10 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node11 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node12 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node13 := merkletree.NewNode([]byte("test1asdfasdf job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node14 := merkletree.NewNode([]byte("test2 job asldkj;fasldkjfasd"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node15 := merkletree.NewNode([]byte("test3 asdfasl;dfasdjob"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node16 := merkletree.NewNode([]byte("tesasdfa;sdasd;laskdjf;alsjflkfj;ast4 job"), &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	tree := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node9, node10, node11, node12, node13, node14, node15, node16, node1, node2, node3, node4, node5, node6, node7, node8})
	return core.NewBlock(*tree, []byte("TestingPreviousHash"), uint64(rand.Int()), difficulty)
}

func (b *BenchmarkEngine) Run() {
	glg.Info("Benchmarking")
	difficulty := 1
	for {
		var avg []float64
		for i := 0; i < 10; i++ {
			start := time.Now()
			block := b.Block(uint8(difficulty))
			end := time.Now()
			block.DeleteFile()
			diff := end.Sub(start)
			avg = append(avg, diff.Seconds())
		}
		var avgSum float64
		for _, val := range avg {
			avgSum += val
		}
		if avgSum/float64(len(avg)) > float64(time.Minute) {
			break
		}
		benchmark := Benchmark{
			AvgTime:    avgSum,
			Difficulty: uint8(difficulty),
		}
		b.AddBenchmark(benchmark)
		difficulty++
	}
	b.SetScore(b.GetData()[len(b.GetData())-1].GetDifficulty())
}

func NewBenchmarkEngine() BenchmarkEngine {
	b := BenchmarkEngine{}
	b.Run()
	return b
}
