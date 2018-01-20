package consensus

import (
	"math/rand"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/core"
)

const Blockrate = 15 // blocks per minute

func Difficulty(benchmarks []benchmark.Benchmark, bc core.BlockChain) int {
	latest := len(bc.GetBlocksWithinMinute())
	for _, val := range benchmarks {
		rate := 60 / val.GetAvgTime()
		if int(rate)+latest <= Blockrate {
			return int(val.GetDifficulty())
		}
	}
	return rand.Intn(int(benchmarks[len(benchmarks)-1].GetDifficulty())) // random difficulty
}
