package benchmark_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/benchmark"
)

func TestNewBenchmark(t *testing.T) {
	b := benchmark.NewBenchmark(46.5, 18)
	assert.NotNil(t, b)
}
