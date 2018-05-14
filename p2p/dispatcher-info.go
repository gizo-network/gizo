package p2p

type DispatcherInfo struct {
	pub        []byte
	neighbours []string
	shut       bool
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

func (w DispatcherInfo) GetShut() bool {
	return w.shut
}

func (w *DispatcherInfo) SetShut(s bool) {
	w.shut = s
}

func (w DispatcherInfo) GetNeighbours() []string {
	return w.neighbours
}

func (w *DispatcherInfo) SetNeighbours(n []string) {
	w.neighbours = n
}

func (w *DispatcherInfo) AddNeighbour(n string) {
	w.neighbours = append(w.neighbours, n)
}
