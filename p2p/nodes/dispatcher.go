package nodes

import (
	"encoding/hex"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gizo-network/gizo/helpers"
	melody "gopkg.in/olahol/melody.v1"

	"github.com/boltdb/bolt"

	"github.com/kpango/glg"

	externalip "github.com/GlenDC/go-external-ip"
	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
)

type Dispatcher struct {
	IP      net.IP
	Port    uint               // port
	Pub     []byte             //public key of the node
	priv    []byte             //private key of the node
	uptime  int64              //time since node has been up
	workers [MaxWorkers]Worker //worker nodes in it's area
	bench   benchmark.Engine   // benchmark of node
	ws      *melody.Melody
	rpc     *rpc.Server
	router  *mux.Router
	jc      *cache.JobCache  //job cache
	bc      *core.BlockChain //blockchain
	db      *bolt.DB         //holds topology table
}

func (d Dispatcher) NodeTypeDispatcher() bool {
	return true
}

func (d Dispatcher) GetIP() net.IP {
	return d.IP
}

func (d Dispatcher) GetPubByte() []byte {
	return d.Pub
}

func (d Dispatcher) GetPubString() string {
	return hex.EncodeToString(d.Pub)
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

func (d Dispatcher) GetPort() int {
	return int(d.Port)
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

func (d Dispatcher) GetWS() *melody.Melody {
	return d.ws
}

func (d *Dispatcher) setWS(m *melody.Melody) {
	d.ws = m
}

func (d Dispatcher) GetRPC() *rpc.Server {
	return d.rpc
}

func (d Dispatcher) setRPC(s *rpc.Server) {
	d.rpc = s
}

func (d Dispatcher) Start() {
	d.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		d.ws.HandleRequest(w, r)
	}).Methods("POST")
	d.router.Handle("/rpc", d.rpc).Methods("POST")
	http.ListenAndServe("localhost"+strconv.FormatInt(int64(d.GetPort()), 10), d.router)
}

func NewDispatcher(port int) *Dispatcher {
	glg.Info("Creating new Dispatcher Node")
	core.InitializeDataPath()

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
			IP:     ip,
			Pub:    pub,
			priv:   priv,
			Port:   uint(port),
			uptime: time.Now().Unix(),
			bench:  bench,
			jc:     jc,
			bc:     bc,
			db:     db,
			router: mux.NewRouter(),
			ws:     melody.New(),
			rpc:    rpc.NewServer(),
		}
	}

	priv, pub = crypt.GenKeys()
	bench = benchmark.NewEngine()
	glg.Warn("Benchmarking done...")
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

	if err != nil {
		glg.Fatal(err)
	}
	return &Dispatcher{
		IP:     ip,
		Pub:    pub,
		priv:   priv,
		Port:   uint(port),
		uptime: time.Now().Unix(),
		bench:  bench,
		jc:     jc,
		bc:     bc,
		db:     db,
		router: mux.NewRouter(),
		ws:     melody.New(),
		rpc:    rpc.NewServer(),
	}
}
