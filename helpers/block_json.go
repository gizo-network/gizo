package helpers

import (
	"encoding/json"

	"github.com/gizo-network/gizo/core"
)

func MarshalBlock(b core.Block) (string, error) {
	temp, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(temp), nil
}
