package core_test

import (
	"testing"

	"github.com/gizo-network/gizo/core"
	"github.com/stretchr/testify/assert"
)

func TestGenesis(t *testing.T) {
	assert.NotNil(t, core.GenesisBlock())
}
