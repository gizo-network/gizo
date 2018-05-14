package p2p

import (
	"encoding/hex"
	"encoding/json"
	"errors"
)

func (d Dispatcher) RpcHttp() {
	d.GetRPC().AddFunction("PeerCount", d.PeerCount)
	d.GetRPC().AddFunction("BlockByHash", d.BlockByHash)
	d.GetRPC().AddFunction("BlockByHeight", d.BlockByHeight)
	d.GetRPC().AddFunction("Latest15Blocks", d.Latest15Blocks)
	d.GetRPC().AddFunction("LatestBlock", d.LatestBlock)
}

func (d Dispatcher) PeerCount() int {
	return len(d.GetPeers())
}

func (d Dispatcher) BlockByHash(hash string) (string, error) {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return "", err
	}
	b, err := d.GetBC().GetBlockInfo(hashBytes)
	if err != nil {
		return "", err
	}
	return string(b.GetBlock().Serialize()), nil
}

func (d Dispatcher) BlockByHeight(height int) (string, error) {
	b, err := d.GetBC().GetBlockByHeight(height)
	if err != nil {
		return "", err
	}
	return string(b.Serialize()), nil
}

func (d Dispatcher) Latest15Blocks() (string, error) {
	blocksBytes, err := json.Marshal(d.GetBC().GetLatest15())
	if err != nil {
		return "", err
	}
	return string(blocksBytes), nil
}

func (d Dispatcher) LatestBlock() string {
	return string(d.GetBC().GetLatestBlock().Serialize())
}

func (d Dispatcher) WriteQueueCount() int {
	return d.GetWriteQ().Size()
}

func (d Dispatcher) Score() float64 {
	return d.GetBench().GetScore()
}

func (d Dispatcher) Peers() []string {
	return d.GetPeersPubs()
}

func (d Dispatcher) PublicKey() string {
	return d.GetPubString()
}

func (d Dispatcher) NewJob(task string, name string, priv bool, privKey string) (string, error) {
	if privKey == "" {
		return "", errors.New("Empty private key")
	}
}
