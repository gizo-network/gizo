package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain()
	bci := bc.iterator()
	assert.NotNil(t, bci.Next())
}
