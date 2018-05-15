package p2p

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Lobarr/lane"
	upnp "github.com/NebulousLabs/go-upnp"
	"github.com/hprose/hprose-golang/rpc"

	"github.com/gizo-network/gizo/core/difficulty"
	"github.com/gizo-network/gizo/core/merkletree"

	"github.com/gizo-network/gizo/job"

	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job/queue"
	funk "github.com/thoas/go-funk"
	melody "gopkg.in/olahol/melody.v1"

	"github.com/boltdb/bolt"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	ErrJobsFull = errors.New("Jobs array full")
)

type Dispatcher struct {
	IP        string
	Port      uint   //port
	Pub       []byte //public key of the node
	priv      []byte //private key of the node
	uptime    int64  //time since node has been up
	jobPQ     *queue.JobPriorityQueue
	workers   map[*melody.Session]*WorkerInfo //worker nodes in dispatcher's area
	peers     map[interface{}]*DispatcherInfo
	workerPQ  *WorkerPriorityQueue
	bench     benchmark.Engine //benchmark of node
	wWS       *melody.Melody   //workers ws server
	dWS       *melody.Melody   //dispatchers ws server
	rpc       *rpc.HTTPService
	router    *mux.Router
	jc        *cache.JobCache  //job cache
	bc        *core.BlockChain //blockchain
	db        *bolt.DB         //holds topology table
	mu        *sync.Mutex
	jobs      []job.Job // holds done jobs and new jobs submitted to the network before being placed in the bc
	interrupt chan os.Signal
	writeQ    *lane.Queue // queue of job (execs) to be written to the db
	centrum   *Centrum
	discover  *upnp.IGD
	new       bool // if true, sends a new dispatcher to centrum else sends a wake with token
}

func (d Dispatcher) GetJobs() []job.Job {
	return d.jobs
}

func (d Dispatcher) watchWriteQ() {
	for {
		if d.GetWriteQ().Empty() == false {
			jobs := d.GetWriteQ().Dequeue()
			d.WriteJobs(jobs.([]job.Job))
		}
	}
}

//WriteJobs writes jobs to the bc
func (d Dispatcher) WriteJobs(jobs []job.Job) {
	nodes := []*merkletree.MerkleNode{}
	for _, job := range jobs {
		nodes = append(nodes, merkletree.NewNode(job, &merkletree.MerkleNode{}, &merkletree.MerkleNode{}))
	}
	block := core.NewBlock(*merkletree.NewMerkleTree(nodes), d.GetBC().GetLatestBlock().GetHeader().GetHash(), d.GetBC().GetNextHeight(), uint8(difficulty.Difficulty(d.GetBenchmarks(), *d.GetBC())), d.GetPubString())
	err := d.GetBC().AddBlock(block)
	if err != nil {
		glg.Fatal(err)
	}
	d.BroadcastPeers(BlockMessage(block.Serialize(), d.GetPrivByte()))
}

func (d *Dispatcher) AddJob(j job.Job) {
	if len(d.GetJobs()) < merkletree.MaxTreeJobs {
		for i, val := range d.GetJobs() {
			if val.GetID() == j.GetID() {
				temp := val
				temp.AddExec(j.GetLatestExec())
				d.jobs[i] = temp
				return
			}
		}
		d.jobs = append(d.jobs, j)
	} else {
		d.GetWriteQ().Enqueue(d.GetJobs())
		d.EmptyJobs()
		d.jobs = append(d.jobs, j)
	}
}

func (d *Dispatcher) EmptyJobs() {
	d.jobs = []job.Job{}
}

func (d Dispatcher) GetWorkerPQ() *WorkerPriorityQueue {
	return d.workerPQ
}

func (d Dispatcher) GetAssignedWorker(hash []byte) *melody.Session {
	for key, val := range d.GetWorkers() {
		if bytes.Compare(val.GetJob().GetJob().GetHash(), hash) == 0 {
			return key
		}
	}
	return nil
}

func (d Dispatcher) NodeTypeDispatcher() bool {
	return true
}

func (d Dispatcher) GetIP() string {
	return d.IP
}

func (d *Dispatcher) SetIP(ip string) {
	d.IP = ip
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

func (d Dispatcher) GetWriteQ() *lane.Queue {
	return d.writeQ
}

func (d Dispatcher) GetJobPQ() *queue.JobPriorityQueue {
	return d.jobPQ
}

func (d Dispatcher) GetWorkers() map[*melody.Session]*WorkerInfo {
	return d.workers
}

func (d Dispatcher) GetWorker(s *melody.Session) *WorkerInfo {
	return d.GetWorkers()[s]
}

func (d *Dispatcher) SetWorker(s *melody.Session, w *WorkerInfo) {
	d.GetWorkers()[s] = w
}

func (d Dispatcher) GetPeers() map[interface{}]*DispatcherInfo {
	return d.peers
}

func (d Dispatcher) GetPeersPubs() []string {
	var temp []string
	for _, info := range d.GetPeers() {
		temp = append(temp, hex.EncodeToString(info.GetPub()))
	}
	return temp
}

func (d Dispatcher) GetPeer(n interface{}) *DispatcherInfo {
	return d.GetPeers()[n]
}

func (d *Dispatcher) AddPeer(s interface{}, n *DispatcherInfo) {
	d.GetPeers()[s] = n
}

func (d Dispatcher) GetBC() *core.BlockChain {
	return d.bc
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

func (d Dispatcher) GetJC() *cache.JobCache {
	return d.jc
}

func (d Dispatcher) GetWWS() *melody.Melody {
	return d.wWS
}

func (d *Dispatcher) setWWS(m *melody.Melody) {
	d.wWS = m
}

func (d Dispatcher) GetDWS() *melody.Melody {
	return d.dWS
}

func (d *Dispatcher) setDWS(m *melody.Melody) {
	d.dWS = m
}

func (d Dispatcher) GetRPC() *rpc.HTTPService {
	return d.rpc
}

func (d Dispatcher) setRPC(s *rpc.HTTPService) {
	d.rpc = s
}

func (d Dispatcher) BroadcastWorkers(m []byte) {
	for s, _ := range d.GetWorkers() {
		s.Write(m)
	}
}

func (d Dispatcher) BroadcastPeers(m []byte) {
	for peer, _ := range d.GetPeers() {
		switch n := peer.(type) {
		case *melody.Session:
			n.Write(m)
			break
		case *websocket.Conn:
			n.WriteMessage(websocket.BinaryMessage, m)
			break
		}
	}
}

func (d Dispatcher) MulticastPeers(m []byte, neigbhours []string) {
	for peer, info := range d.GetPeers() {
		if funk.ContainsString(neigbhours, hex.EncodeToString(info.GetPub())) {
			switch n := peer.(type) {
			case *melody.Session:
				n.Write(m)
				break
			case *websocket.Conn:
				n.WriteMessage(websocket.BinaryMessage, m)
				break
			}
		}
	}
}

func (d Dispatcher) WorkerExists(s *melody.Session) bool {
	_, ok := d.workers[s]
	return ok
}

func (d *Dispatcher) deployJobs() {
	for {
		if d.GetWorkerPQ().getPQ().Empty() == false {
			if d.GetJobPQ().GetPQ().Empty() == false {
				d.mu.Lock()
				w := d.GetWorkerPQ().Pop()
				if !d.GetWorker(w).GetShut() {
					j := d.GetJobPQ().Pop()
					if j.GetExec().GetStatus() != job.CANCELLED {
						j.GetExec().SetBy(d.GetWorker(w).GetPub())
						d.GetWorker(w).Assign(&j)
						glg.Info("P2P: dispatched job")
						w.Write(JobMessage(j.Serialize(), d.GetPrivByte()))
					} else {
						j.ResultsChan() <- j
					}
				} else {
					delete(d.GetWorkers(), w)
				}
				d.mu.Unlock()
			}
		}
	}
}

func (d Dispatcher) wPeerTalk() {
	d.wWS.HandleDisconnect(func(s *melody.Session) {
		d.mu.Lock()
		glg.Info("Dispatcher: worker disconnected")
		if d.GetWorker(s).GetJob() != nil {
			d.GetJobPQ().PushItem(*d.GetWorker(s).GetJob(), job.HIGH)
		}
		d.GetWorker(s).SetShut(true)
		d.mu.Unlock()
	})
	d.wWS.HandleMessageBinary(func(s *melody.Session, message []byte) {
		m := DeserializePeerMessage(message)
		switch m.GetMessage() {
		case HELLO:
			d.mu.Lock()
			if len(d.GetWorkers()) < MaxWorkers {
				glg.Info("Dispatcher: worker connected")
				d.SetWorker(s, NewWorkerInfo(hex.EncodeToString(m.GetPayload())))
				s.Write(HelloMessage(d.GetPubByte()))
				d.centrum.ConnectWorker()
				d.GetWorkerPQ().Push(s, 0)
			} else {
				s.Write(ConnFullMessage())
			}
			d.mu.Unlock()
			break
		case RESULT:
			d.mu.Lock()
			if m.VerifySignature(d.GetWorker(s).GetPub()) {
				glg.Info("P2P: received result")
				exec := job.DeserializeExec(m.GetPayload())
				d.GetWorker(s).GetJob().SetExec(&exec)
				d.GetWorker(s).GetJob().ResultsChan() <- *d.GetWorker(s).GetJob()
				j := d.GetWorker(s).GetJob().GetJob()
				d.GetWorker(s).SetJob(nil)
				j.AddExec(exec)
				//TODO: send to requester
				d.AddJob(j)
			} else {
				d.GetJobPQ().PushItem(*d.GetWorker(s).GetJob(), job.HIGH)
			}
			if !d.GetWorker(s).GetShut() {
				d.GetWorkerPQ().Push(s, 0)
			}
			d.mu.Unlock()
			break
		case SHUT:
			d.mu.Lock()
			d.GetWorker(s).SetShut(true)
			s.Write(ShutAckMessage(d.GetPrivByte()))
			d.centrum.DisconnectWorker()
			d.mu.Unlock()
			break
		default:
			s.Write(InvalidMessage())
			break
		}
	})
}

func (d Dispatcher) dPeerTalk() {
	d.dWS.HandleDisconnect(func(s *melody.Session) {
		d.mu.Lock()
		info := d.GetPeer(s)
		if info != nil {
			glg.Info("Dispatcher: peer disconnected")
			d.BroadcastPeers(PeerDisconnectMessage(info.GetPub(), d.GetPrivByte()))
			delete(d.GetPeers(), s)
		}
		d.mu.Unlock()
	})
	d.dWS.HandleMessageBinary(func(s *melody.Session, message []byte) {
		m := DeserializePeerMessage(message)
		switch m.GetMessage() {
		case HELLO:
			d.mu.Lock()
			info := DeserializeDispatcherHello(m.GetPayload())
			d.AddPeer(s, &DispatcherInfo{pub: info.GetPub(), peers: info.GetPeers()})
			s.Write(HelloMessage(NewDispatcherHello(d.GetPubByte(), d.GetPeersPubs()).Serialize()))
			d.mu.Unlock()
			break
		case BLOCK:
			d.mu.Lock()
			if m.VerifySignature(hex.EncodeToString(d.GetPeer(s).GetPub())) {
				b, err := core.DeserializeBlock(m.GetPayload())
				if err != nil {
					glg.Fatal(err)
				}
				err = b.Export()
				if err != nil {
					glg.Fatal(err)
				}
				d.GetBC().AddBlock(b)
				var peerToRecv []string
				for _, info := range d.GetPeers() {
					//! avoids broadcast storms by not sending block back to sender and to neigbhours that are not directly connected to sender
					if !funk.ContainsString(info.GetPeers(), hex.EncodeToString(d.GetPeer(s).GetPub())) && bytes.Compare(info.GetPub(), d.GetPeer(s).GetPub()) != 0 {
						peerToRecv = append(peerToRecv, hex.EncodeToString(info.GetPub()))
					}
				}
				d.MulticastPeers(BlockMessage(m.GetPayload(), d.GetPrivByte()), peerToRecv)
			}
			d.mu.Unlock()
			break
		case BLOCKREQ:
			d.mu.Lock()
			if m.VerifySignature(hex.EncodeToString(d.GetPeer(s).GetPub())) {
				blockinfo, _ := d.GetBC().GetBlockInfo(m.GetPayload())
				s.Write(BlockResMessage(blockinfo.GetBlock().Serialize(), d.GetPrivByte()))
			}
			d.mu.Unlock()
			break
		case PEERCONNECT:
			d.mu.Lock()
			if m.VerifySignature(hex.EncodeToString(d.GetPeer(s).GetPub())) {
				d.GetPeer(s).AddPeer(hex.EncodeToString(m.GetPayload()))
			}
			d.mu.Unlock()
			break
		case PEERDISCONNECT:
			d.mu.Lock()
			if m.VerifySignature(hex.EncodeToString(d.GetPeer(s).GetPub())) {
				peers := d.GetPeer(s).GetPeers()
				for i, peer := range peers {
					if peer == hex.EncodeToString(m.GetPayload()) {
						d.GetPeer(s).SetPeers(append(peers[:i], peers[i+1:]...))
						break
					}
				}
			}
			d.mu.Unlock()
			break
		default:
			s.Write(InvalidMessage())
			break
		}
	})
}

func (d Dispatcher) HandleNodeConnect(conn *websocket.Conn) {
	conn.WriteMessage(websocket.BinaryMessage, HelloMessage(NewDispatcherHello(d.GetPubByte(), d.GetPeersPubs()).Serialize()))
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			//TODO: handle syncer disconnect - use next best version
			d.mu.Lock()
			glg.Info("Dispatcher: peer disconnected")
			info := d.GetPeer(conn)
			d.BroadcastPeers(PeerDisconnectMessage(info.GetPub(), d.GetPrivByte()))
			delete(d.GetPeers(), conn)
			d.mu.Unlock()
		}
		m := DeserializePeerMessage(message)
		switch m.GetMessage() {
		case HELLO:
			d.mu.Lock()
			peerInfo := DeserializeDispatcherHello(m.GetPayload())
			if bytes.Compare(d.GetPeer(conn).GetPub(), peerInfo.GetPub()) == 0 {
				d.GetPeer(conn).SetPeers(peerInfo.GetPeers())
			} else {
				delete(d.GetPeers(), conn)
				conn.Close()
			}
			d.mu.Unlock()
			break
		case BLOCK:
			d.mu.Lock()
			if m.VerifySignature(hex.EncodeToString(d.GetPeer(conn).GetPub())) {
				b, err := core.DeserializeBlock(m.GetPayload())
				if err != nil {
					glg.Fatal(err)
				}
				err = b.Export()
				if err != nil {
					glg.Fatal(err)
				}
				d.GetBC().AddBlock(b)
				var peerToRecv []string
				for _, info := range d.GetPeers() {
					//! avoids broadcast storms by not sending block back to sender and to neigbhours that are not directly connected to sender
					if !funk.ContainsString(info.GetPeers(), hex.EncodeToString(d.GetPeer(conn).GetPub())) && bytes.Compare(info.GetPub(), d.GetPeer(conn).GetPub()) != 0 {
						peerToRecv = append(peerToRecv, hex.EncodeToString(info.GetPub()))
					}
				}
				d.MulticastPeers(BlockMessage(m.GetPayload(), d.GetPrivByte()), peerToRecv)
			}
			d.mu.Unlock()
			break
		case BLOCKRES:
			if m.VerifySignature(hex.EncodeToString(d.GetPeer(conn).GetPub())) {
				b, err := core.DeserializeBlock(m.GetPayload())
				if err != nil {
					glg.Fatal(err)
				}
				err = b.Export()
				if err != nil {
					glg.Fatal(err)
				}
				d.GetBC().AddBlock(b)
			}
			break
		case PEERCONNECT:
			d.mu.Lock()
			if m.VerifySignature(hex.EncodeToString(d.GetPeer(conn).GetPub())) {
				d.GetPeer(conn).AddPeer(hex.EncodeToString(m.GetPayload()))
			}
			d.mu.Unlock()
			break
		case PEERDISCONNECT:
			d.mu.Lock()
			if m.VerifySignature(hex.EncodeToString(d.GetPeer(conn).GetPub())) {
				peers := d.GetPeer(conn).GetPeers()
				for i, peer := range peers {
					if peer == hex.EncodeToString(m.GetPayload()) {
						d.GetPeer(conn).SetPeers(append(peers[:i], peers[i+1:]...))
					}
				}
			}
			d.mu.Unlock()
			break
		default:
			conn.WriteMessage(websocket.BinaryMessage, InvalidMessage())
			break
		}
	}
}

func (d Dispatcher) WatchInterrupt() {
	select {
	case i := <-d.interrupt:
		glg.Warn("Dispatcher: interrupt detected")
		switch i {
		case syscall.SIGINT, syscall.SIGTERM:
			//TODO: get jobs from workers, add all jobs held in memory to block, broadcast block to the network
			res, err := d.centrum.Sleep()
			if err != nil {
				glg.Fatal(err)
			} else if res["status"].(string) != "success" {
				glg.Fatal("Centrum: " + res["status"].(string))
			}
			d.BroadcastWorkers(ShutMessage(d.GetPrivByte()))
			time.Sleep(time.Second * 3) // give neighbors and workers 3 seconds to disconnect
			os.Exit(0)
		case syscall.SIGQUIT:
			os.Exit(1)
		}
	}
}

func (d Dispatcher) Start() {
	if !d.GetBC().Verify() {
		glg.Fatal("Dispatcher: blockchain not verified")
	}
	go d.deployJobs()
	go d.watchWriteQ()
	go d.WatchInterrupt()
	d.GetDispatchersAndSync()
	d.wWS.Upgrader.ReadBufferSize = 100000
	d.wWS.Upgrader.WriteBufferSize = 100000
	d.wWS.Config.MessageBufferSize = 100000
	d.wWS.Config.MaxMessageSize = 100000
	d.wWS.Upgrader.EnableCompression = true
	d.dWS.Upgrader.ReadBufferSize = 100000
	d.dWS.Upgrader.WriteBufferSize = 100000
	d.dWS.Config.MessageBufferSize = 100000
	d.dWS.Config.MaxMessageSize = 100000
	d.dWS.Upgrader.EnableCompression = true
	d.router.HandleFunc("/d", func(w http.ResponseWriter, r *http.Request) {
		d.dWS.HandleRequest(w, r)
	})
	d.router.HandleFunc("/w", func(w http.ResponseWriter, r *http.Request) {
		d.wWS.HandleRequest(w, r)
	})
	d.wPeerTalk()
	d.dPeerTalk()
	d.RpcHttp()
	d.router.Handle("/rpc", d.GetRPC()).Methods("POST")
	status := make(map[string]string)
	status["status"] = "running"
	status["pub"] = d.GetPubString()
	statusBytes, err := json.Marshal(status)
	if err != nil {
		glg.Fatal(err)
	}
	d.router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write(statusBytes)
	})
	d.router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write(NewVersion(GizoVersion, int(d.GetBC().GetLatestHeight()), d.GetBC().GetBlockHashesHex()).Serialize())
	})

	err = d.discover.Forward(uint16(d.GetPort()), "gizo dispatcher node")
	if err != nil {
		log.Fatal(err)
	}
	if d.new {
		go d.Register()
	} else {
		res, err := d.centrum.Wake()
		if err != nil {
			glg.Fatal(err)
		} else if res["status"].(string) != "success" {
			glg.Fatal("Centrum: " + res["status"].(string))
		}
	}
	fmt.Println(http.ListenAndServe(":"+strconv.FormatInt(int64(d.GetPort()), 10), d.router))
}

func (d Dispatcher) SaveToken() {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(NodeBucket))
		if err := b.Put([]byte("token"), []byte(d.centrum.GetToken())); err != nil {
			glg.Fatal(err)
		}
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
}

func (d Dispatcher) Register() {
	time.Sleep(time.Second * 1)
	if err := d.centrum.NewDisptcher(d.GetPubString(), d.GetIP(), int(d.GetPort())); err != nil {
		glg.Warn("Centrum: unable to get on network")
		glg.Fatal("Centrum: " + err.Error())
	}
	d.SaveToken()
}

func (d *Dispatcher) GetDispatchersAndSync() {
	time.Sleep(time.Second * 2)
	res := d.centrum.GetDispatchers()
	syncVersion := new(Version)
	syncPeer := new(websocket.Conn)
	dispatchers, ok := res["dispatchers"]
	if !ok {
		glg.Warn(ErrNoDispatchers)
		return
	}
	for _, dispatcher := range dispatchers.([]string) {
		addr, err := ParseAddr(dispatcher)
		if err == nil && addr["pub"].(string) != d.GetPubString() {
			var v Version
			wsURL := fmt.Sprintf("ws://%v:%v/d", addr["ip"], addr["port"])
			versionURL := fmt.Sprintf("http://%v:%v/rpc", addr["ip"], addr["port"])
			dailer := websocket.Dialer{
				Proxy:           http.ProxyFromEnvironment,
				ReadBufferSize:  10000,
				WriteBufferSize: 10000,
			}
			conn, _, err := dailer.Dial(wsURL, nil)
			if err != nil {
				glg.Fatal(err)
			}
			conn.EnableWriteCompression(true)
			pubBytes, err := hex.DecodeString(addr["pub"].(string))
			if err != nil {
				glg.Fatal(err)
			}
			d.AddPeer(conn, NewDispatcherInfo(pubBytes))
			go d.HandleNodeConnect(conn)
			_, err = s.New().Get(versionURL).ReceiveSuccess(&v)
			if err != nil {
				glg.Fatal(err)
			}
			if syncVersion.GetHeight() < v.GetHeight() {
				syncVersion = &v
				syncPeer = conn
			}
		}
	}
	if syncVersion.GetHeight() != 0 {
		glg.Warn("Dispatcher: node sync in progress")
		blocks := d.GetBC().GetBlockHashesHex()
		for _, hash := range syncVersion.GetBlocks() {
			if !funk.ContainsString(blocks, hash) {
				hashBytes, err := hex.DecodeString(hash)
				if err != nil {
					glg.Fatal(err)
				}
				syncPeer.WriteMessage(websocket.BinaryMessage, BlockReqMessage(hashBytes, d.GetPrivByte()))
			}
		}
	}
}

func NewDispatcher(port int) *Dispatcher {
	glg.Info("Creating Dispatcher Node")
	core.InitializeDataPath()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	var bench benchmark.Engine
	var priv, pub []byte
	var token string
	discover, err := upnp.Discover()
	if err != nil {
		glg.Fatal(err)
	}

	ip, err := discover.ExternalIP()
	if err != nil {
		log.Fatal(err)
	}

	centrum := NewCentrum()

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
			token = string(b.Get([]byte("token")))
			return nil
		})
		if err != nil {
			glg.Fatal(err)
		}
		centrum.SetToken(token)
		bc := core.CreateBlockChain(hex.EncodeToString(pub))
		jc := cache.NewJobCache(bc)
		return &Dispatcher{
			IP:        ip,
			Pub:       pub,
			priv:      priv,
			Port:      uint(port),
			uptime:    time.Now().Unix(),
			bench:     bench,
			jobPQ:     queue.NewJobPriorityQueue(),
			workers:   make(map[*melody.Session]*WorkerInfo),
			workerPQ:  NewWorkerPriorityQueue(),
			peers:     make(map[interface{}]*DispatcherInfo),
			jc:        jc,
			bc:        bc,
			db:        db,
			router:    mux.NewRouter(),
			wWS:       melody.New(),
			dWS:       melody.New(),
			rpc:       rpc.NewHTTPService(),
			mu:        new(sync.Mutex),
			interrupt: interrupt,
			writeQ:    lane.NewQueue(),
			centrum:   centrum,
			discover:  discover,
			new:       false,
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

	if err != nil {
		glg.Fatal(err)
	}
	bc := core.CreateBlockChain(hex.EncodeToString(pub))
	jc := cache.NewJobCache(bc)
	return &Dispatcher{
		IP:        ip,
		Pub:       pub,
		priv:      priv,
		Port:      uint(port),
		uptime:    time.Now().Unix(),
		bench:     bench,
		jobPQ:     queue.NewJobPriorityQueue(),
		workers:   make(map[*melody.Session]*WorkerInfo),
		workerPQ:  NewWorkerPriorityQueue(),
		peers:     make(map[interface{}]*DispatcherInfo),
		jc:        jc,
		bc:        bc,
		db:        db,
		router:    mux.NewRouter(),
		wWS:       melody.New(),
		dWS:       melody.New(),
		rpc:       rpc.NewHTTPService(),
		mu:        new(sync.Mutex),
		interrupt: interrupt,
		writeQ:    lane.NewQueue(),
		centrum:   centrum,
		discover:  discover,
		new:       true,
	}
}
