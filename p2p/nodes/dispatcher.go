package nodes

import (
	"encoding/hex"
	"net"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/boltdb/bolt"

	"github.com/kpango/glg"

	externalip "github.com/GlenDC/go-external-ip"
	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/crypt"
)

const (
	MaxWorkers     = 128
	DefaultWSPort  = 9998
	DefaultRPCPort = 9999
)

type Dispatcher struct {
	ip      net.IP
	pub     []byte             //public key of the node
	priv    []byte             //private key of the node
	workers [MaxWorkers]Worker //worker nodes in it's area
	rpc     uint               // rpc port
	ws      uint               // ws port
	uptime  int64              //time since node has been up
	bench   benchmark.Engine   // benchmark of node
	jc      *cache.JobCache    //job cache
	bc      *core.BlockChain   //blockchain
	db      *bolt.DB           //holds topology table
}

func (d Dispatcher) GetIP() net.IP {
	return d.ip
}

func (d Dispatcher) GetPubByte() []byte {
	return d.pub
}

func (d Dispatcher) GetPubString() string {
	return hex.EncodeToString(d.pub)
}

func (d Dispatcher) GetPrivByte() []byte {
	return d.priv
}

func (d Dispatcher) GetPrivString() string {
	return hex.EncodeToString(d.priv)
}

func (d Dispatcher) GetWorkers() [MaxWorkers]Worker {
	return d.workers
}

func (d Dispatcher) GetRPCPort() int {
	return int(d.rpc)
}

func (d Dispatcher) GetWSPort() int {
	return int(d.ws)
}

func (d Dispatcher) GetUptme() int64 {
	return d.uptime
}

func (d Dispatcher) GetUptimeString() string {
	return time.Unix(d.uptime, 0).Sub(time.Now()).String()
}

func (d Dispatcher) GetBench() benchmark.Engine {
	return d.bench
}

func (d Dispatcher) GetBenchmarks() []benchmark.Benchmark {
	return d.bench.GetData()
}

func NewDispatcher(rpc, ws int) *Dispatcher {
	core.InitializeDataPath()
	if reflect.ValueOf(rpc).IsNil() {
		rpc = DefaultRPCPort
	}

	if reflect.ValueOf(ws).IsNil() {
		ws = DefaultWSPort
	}

	var bench benchmark.Engine
	var priv, pub []byte
	bc := core.CreateBlockChain(hex.EncodeToString(pub))
	jc := cache.NewJobCache(bc)
	ip, err := externalip.DefaultConsensus(nil, nil).ExternalIP()
	if err != nil {
		glg.Fatal(err)
	}

	var dbFile string
	if os.Getenv("ENV") == "dev" {
		dbFile = path.Join(core.IndexPathDev, NodeDB)
	} else {
		dbFile = path.Join(core.IndexPathProd, NodeDB)
	}

	if helpers.FileExists(dbFile) {
		glg.Warn("Dispatcher: using existing keypair and benchmark")
		db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
		if err != nil {
			glg.Fatal(err)
		}
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(NodeBucket))
			priv = b.Get([]byte("priv"))
			pub = b.Get([]byte("pub"))
			bench = benchmark.DeserializeBenchmarkEngine(b.Get([]byte("benchmark")))
			return nil
		})
		if err != nil {
			glg.Fatal(err)
		}
		return &Dispatcher{
			ip:     ip,
			pub:    pub,
			priv:   priv,
			rpc:    uint(rpc),
			ws:     uint(ws),
			uptime: time.Now().Unix(),
			bench:  bench,
			jc:     jc,
			bc:     bc,
			db:     db,
		}
	}

	priv, pub = crypt.GenKeys()
	bench = benchmark.NewEngine()

	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
	if err != nil {
		glg.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(NodeBucket))
		if err != nil {
			glg.Fatal(err)
		}

		if err = b.Put([]byte("benchmark"), bench.Serialize()); err != nil {
			glg.Fatal(err)
		}

		if err = b.Put([]byte("priv"), priv); err != nil {
			glg.Fatal(err)
		}

		if err = b.Put([]byte("pub"), pub); err != nil {
			glg.Fatal(err)
		}
		return nil
	})

	return &Dispatcher{
		ip:     ip,
		pub:    pub,
		priv:   priv,
		rpc:    uint(rpc),
		ws:     uint(ws),
		uptime: time.Now().Unix(),
		bench:  bench,
		jc:     jc,
		bc:     bc,
		db:     db,
	}
}
