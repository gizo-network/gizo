package p2p

import "errors"

const (
	NodeDB           = "nodeinfo.db"
	NodeBucket       = "node"
	DispatcherScheme = "gizo" //FIXME: use better one
	MaxWorkers       = 128
	DefaultPort      = 9999
	CentrumURL       = "https://61a1e577.ngrok.io"
	GizoVersion      = 1
)

// node states
const (
	// when a node is not connected to the network
	DOWN = "DOWN"
	// worker - when a worker connects to a dispatchers standard area
	// dispatcher - when an adjacency is created and topology table, peer table and blockchain have not been synced
	INIT = "INIT"
	// worker - when a node starts receiving and crunching jobs
	LIVE = "LIVE"
	// dispatcher - when an adjacency is created and topology table, peer table and blockchain have been sync
	FULL = "FULL"
)

var (
	ErrNoDispatchers = errors.New("Centrum: no dispatchers available")
)
