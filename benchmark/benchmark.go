package benchmark

//Benchmark holds a difficulty and the average time the node takes to run the difficulty
type Benchmark struct {
	AvgTime    float64
	Difficulty uint8
}

// NewBenchmark returns a new benchmark
func NewBenchmark(avgTime float64, difficulty uint8) Benchmark {
	return Benchmark{
		AvgTime:    avgTime,
		Difficulty: difficulty,
	}
}

//GetAvgTime returns the avgTime
func (b Benchmark) GetAvgTime() float64 {
	return b.AvgTime
}

//SetAvgTime set's the avgTime
func (b *Benchmark) SetAvgTime(avg float64) {
	b.AvgTime = avg
}

//GetDifficulty returns the difficulty
func (b Benchmark) GetDifficulty() uint8 {
	return b.Difficulty
}

//SetDifficulty set's the difficulty
func (b *Benchmark) SetDifficulty(d uint8) {
	b.Difficulty = d
}
