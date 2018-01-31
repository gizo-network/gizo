package benchmark

//Benchmark holds a difficulty and the average time the node takes to run the difficulty
type Benchmark struct {
	AvgTime    float64
	Difficulty uint8
}

func NewBenchmark(avgTime float64, difficulty uint8) Benchmark {
	return Benchmark{
		AvgTime:    avgTime,
		Difficulty: difficulty,
	}
}

func (b Benchmark) GetAvgTime() float64 {
	return b.AvgTime
}

func (b *Benchmark) SetAvgTime(avg float64) {
	b.AvgTime = avg
}

func (b Benchmark) GetDifficulty() uint8 {
	return b.Difficulty
}

func (b *Benchmark) SetDifficulty(d uint8) {
	b.Difficulty = d
}
