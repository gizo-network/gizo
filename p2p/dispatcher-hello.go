package p2p

import (
	"encoding/json"

	"github.com/kpango/glg"
)

type DispatcherHello struct {
	Pub   []byte
	Peers []string
}

func NewDispatcherHello(pub []byte, p []string) DispatcherHello {
	return DispatcherHello{Pub: pub, Peers: p}
}

func (d DispatcherHello) GetPub() []byte {
	return d.Pub
}

func (d DispatcherHello) GetPeers() []string {
	return d.Peers
}

func (d *DispatcherHello) SetPub(pub []byte) {
	d.Pub = pub
}

func (d *DispatcherHello) SetPeers(n []string) {
	d.Peers = n
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
