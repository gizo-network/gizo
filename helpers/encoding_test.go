package helpers_test

import (
	"testing"

	"github.com/gizo-network/gizo/helpers"
	"github.com/stretchr/testify/assert"
)

func TestEncode64(t *testing.T) {
	assert.NotNil(t, helpers.Encode64([]byte("testing")))
}

func TestDecode64(t *testing.T) {
	b := []byte("testing")
	enc := helpers.Encode64(b)
	assert.Equal(t, b, helpers.Decode64(enc))
}
