package p2p

type DispatcherInfo struct {
	pub  string
	shut bool
}

func NewDispatcherInfo(pub string) *DispatcherInfo {
	return &DispatcherInfo{pub: pub}
}

func (w DispatcherInfo) GetPub() string {
	return w.pub
}

func (w *DispatcherInfo) SetPub(pub string) {
	w.pub = pub
}

func (w DispatcherInfo) GetShut() bool {
	return w.shut
}

func (w *DispatcherInfo) SetShut(s bool) {
	w.shut = s
}
