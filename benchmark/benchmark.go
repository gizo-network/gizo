package benchmark

//Benchmark holds a difficulty and the average time the node takes to run the difficulty
type Benchmark struct {
	avgTime    float64
	difficulty uint8
}

func NewBenchmark(avgTime float64, difficulty uint8) Benchmark {
	return Benchmark{
		avgTime:    avgTime,
		difficulty: difficulty,
	}
}

func (b Benchmark) GetAvgTime() float64 {
	return b.avgTime
}

func (b *Benchmark) SetAvgTime(avg float64) {
	b.avgTime = avg
}

func (b Benchmark) GetDifficulty() uint8 {
	return b.difficulty
}

func (b *Benchmark) SetDifficulty(d uint8) {
	b.difficulty = d
}
