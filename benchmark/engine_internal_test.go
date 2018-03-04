package benchmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	b := Engine{}
	assert.NotNil(t, b.block(10))
}
