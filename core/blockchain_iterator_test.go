package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	RemoveDataPath()
	bc := CreateBlockChain()
	bci := bc.iterator()
	assert.NotNil(t, bci.Next())
}
