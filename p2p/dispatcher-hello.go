package p2p

import (
	"encoding/json"

	"github.com/kpango/glg"
)

type DispatcherHello struct {
	Pub        []byte
	Neighbours []string
}

func NewDispatcherHello(pub []byte, n []string) DispatcherHello {
	return DispatcherHello{Pub: pub, Neighbours: n}
}

func (d DispatcherHello) GetPub() []byte {
	return d.Pub
}

func (d DispatcherHello) GetNeighbours() []string {
	return d.Neighbours
}

func (d *DispatcherHello) SetPub(pub []byte) {
	d.Pub = pub
}

func (d *DispatcherHello) SetNeighbours(n []string) {
	d.Neighbours = n
}

func (d DispatcherHello) Serialize() []byte {
	bytes, err := json.Marshal(d)
	if err != nil {
		glg.Fatal(err)
	}
	return bytes
}

func DeserializeDispatcherHello(b []byte) DispatcherHello {
	var temp DispatcherHello
	err := json.Unmarshal(b, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}
