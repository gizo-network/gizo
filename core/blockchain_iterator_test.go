package core

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	godotenv.Load()
	RemoveDataPath()
	bc := CreateBlockChain()
	bci := bc.iterator()
	assert.NotNil(t, bci.Next())
}
