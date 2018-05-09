package p2p

import (
	"encoding/json"

	"github.com/kpango/glg"
)

type Version struct {
	Version int
	Height  int
	Blocks  []string
}

func NewVersion(version int, height int, blocks []string) Version {
	return Version{
		Version: version,
		Height:  height,
		Blocks:  blocks,
	}
}

func (v Version) GetVersion() int {
	return v.Version
}

func (v Version) GetHeight() int {
	return v.Height
}

func (v Version) GetBlocks() []string {
	return v.Blocks
}

func (v Version) Serialize() []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		glg.Fatal(err)
	}
	return bytes
}

func DeserializeVersion(b []byte) Version {
	var temp Version
	err := json.Unmarshal(b, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}
