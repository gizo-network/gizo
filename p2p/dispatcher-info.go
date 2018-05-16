package p2p

type DispatcherInfo struct {
	pub   []byte
	peers []string
}

func NewDispatcherInfo(pub []byte) *DispatcherInfo {
	return &DispatcherInfo{pub: pub}
}

func (w DispatcherInfo) GetPub() []byte {
	return w.pub
}

func (w *DispatcherInfo) SetPub(pub []byte) {
	w.pub = pub
}

func (w DispatcherInfo) GetPeers() []string {
	return w.peers
}

func (w *DispatcherInfo) SetPeers(n []string) {
	w.peers = n
}

func (w *DispatcherInfo) AddPeer(n string) {
	w.peers = append(w.peers, n)
}
