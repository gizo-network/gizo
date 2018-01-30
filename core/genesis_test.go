package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenesis(t *testing.T) {
	assert.NotNil(t, GenesisBlock())
}
