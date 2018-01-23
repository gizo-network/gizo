package consensus

import (
	"math/rand"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/core"
)

const Blockrate = 15 // blocks per minute

//Difficulty returns a difficulty based on the blockrate and the number of blocks in the last minute
func Difficulty(benchmarks []benchmark.Benchmark, bc core.BlockChain) int {
	latest := len(bc.GetBlocksWithinMinute())
	for _, val := range benchmarks {
		rate := 60 / val.GetAvgTime()
		if int(rate)+latest <= Blockrate {
			return int(val.GetDifficulty())
		}
	}
	return rand.Intn(int(benchmarks[len(benchmarks)-1].GetDifficulty())) // returns random difficulty if difficulty can't be determined
}
