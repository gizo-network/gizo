package benchmark

import (
	"math/rand"
	"sync"
	"time"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
)

type BenchmarkEngine struct {
	Data  []Benchmark
	Score float64
	mu    sync.Mutex
}

func (b *BenchmarkEngine) SetScore(s float64) {
	b.Score = s
}

func (b BenchmarkEngine) GetScore() float64 {
	return b.Score
}

func (b *BenchmarkEngine) AddBenchmark(benchmark Benchmark) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Data = append(b.Data, benchmark)
}

func (b BenchmarkEngine) GetData() []Benchmark {
	b.mu.Lock()
	defer b.mu.Unlock()
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

// Run spins up the benchmark engine
func (b *BenchmarkEngine) Run() {
	glg.Info("Benchmarking node")
	done := false
	var wg sync.WaitGroup
	difficulty := 10 //! difficulty starts at 10
	for done == false {
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(myDifficulty int) {
				var avg []float64
				var mu sync.Mutex
				var mineWG sync.WaitGroup
				for j := 0; j < 5; j++ {
					mineWG.Add(1)
					go func() {
						start := time.Now()
						block := b.Block(uint8(myDifficulty))
						end := time.Now()
						block.DeleteFile()
						diff := end.Sub(start).Seconds()
						mu.Lock()
						avg = append(avg, diff)
						mu.Unlock()
						mineWG.Done()
					}()
				}
				mineWG.Wait()
				var avgSum float64
				for _, val := range avg {
					avgSum += val
				}
				average := avgSum / float64(len(avg))
				if average > 60 {
					done = true
				} else {
					benchmark := Benchmark{
						AvgTime:    average,
						Difficulty: uint8(myDifficulty),
					}
					b.AddBenchmark(benchmark)
				}
				wg.Done()
			}(difficulty)
			difficulty++
		}
		wg.Wait()
	}
	score := float64(b.GetData()[len(b.GetData())-1].GetDifficulty()) - 10  //! 10 is subtracted to allow the score start from 1 since difficulty starts at 10
	scoreDecimal := 1 - b.GetData()[len(b.GetData())-1].GetAvgTime()/100 // determine decimal part of score 
	b.SetScore(score + scoreDecimal)

}

func NewBenchmarkEngine() BenchmarkEngine {
	b := BenchmarkEngine{}
	b.Run()
	return b
}
