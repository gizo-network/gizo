package nodes

import (
	"net"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
)

type Dispatcher struct {
	IP     net.IP
	ID     []byte            //public key of the node
	Area   []Worker          //worker nodes in it's area
	RPC    uint              // rpc port
	WS     uint              // ws port
	UpTime int64             //time since node has been up
	bench  *benchmark.Engine // benchmark of node
	jc     *cache.JobCache   //job cache
	bc     *core.BlockChain  //blockchain
}
