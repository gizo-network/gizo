package p2p

import (
	"encoding/json"

	"github.com/kpango/glg"
)

//DispatcherInfo message used to create and maintain adjacency between dispatcher nodes
type DispatcherInfo struct {
	Pub   []byte
	Peers []string
}

func NewDispatcherInfo(pub []byte, p []string) *DispatcherInfo {
	return &DispatcherInfo{Pub: pub, Peers: p}
}

func (d DispatcherInfo) GetPub() []byte {
	return d.Pub
}

func (d DispatcherInfo) GetPeers() []string {
	return d.Peers
}

func (d *DispatcherInfo) SetPub(pub []byte) {
	d.Pub = pub
}

func (d *DispatcherInfo) SetPeers(n []string) {
	d.Peers = n
}

func (w *DispatcherInfo) AddPeer(n string) {
	w.Peers = append(w.Peers, n)
}

func (d DispatcherInfo) Serialize() []byte {
	bytes, err := json.Marshal(d)
	if err != nil {
		glg.Fatal(err)
	}
	return bytes
}

func DeserializeDispatcherInfo(b []byte) DispatcherInfo {
	var temp DispatcherInfo
	err := json.Unmarshal(b, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}
