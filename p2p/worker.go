package p2p

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job/queue/qItem"
	"github.com/gorilla/websocket"
	"github.com/kpango/glg"
)

type Worker struct {
	Pub        []byte //public key of the node
	Dispatcher string
	shortlist  []string // array of dispatchers received from centrum
	priv       []byte   //private key of the node
	uptime     int64    //time since node has been up
	conn       *websocket.Conn
	interrupt  chan os.Signal
	shutdown   chan struct{}
	busy       bool
	state      string
	item       qItem.Item
}

func (w Worker) GetItem() qItem.Item {
	return w.item
}

func (w *Worker) SetItem(i qItem.Item) {
	w.item = i
}

func (w Worker) GetShortlist() []string {
	return w.shortlist
}

func (w *Worker) SetShortlist(s []string) {
	w.shortlist = s
}

func (w Worker) GetBusy() bool {
	return w.busy
}

func (w *Worker) SetBusy(b bool) {
	w.busy = b
}

func (w Worker) NodeTypeDispatcher() bool {
	return false
}

func (w Worker) GetState() string {
	return w.state
}

func (w *Worker) SetState(s string) {
	w.state = s
}

func (w Worker) GetPubByte() []byte {
	return w.Pub
}

func (w Worker) GetPubString() string {
	return hex.EncodeToString(w.Pub)
}

func (w Worker) GetPrivByte() []byte {
	return w.priv
}

func (w Worker) GetPrivString() string {
	return hex.EncodeToString(w.priv)
}

func (w Worker) GetUptme() int64 {
	return w.uptime
}

func (w Worker) GetDispatcher() string {
	return w.Dispatcher
}

func (w *Worker) SetDispatcher(d string) {
	w.Dispatcher = d
}

func (w Worker) GetUptimeString() string {
	return time.Unix(w.uptime, 0).Sub(time.Now()).String()
}

func (w *Worker) Start() {
	//TODO: implemented cancel and getstatus
	w.GetDispatchers()
	w.Connect()
	go w.WatchInterrupt()
	w.conn.WriteMessage(websocket.BinaryMessage, HelloMessage(w.GetPubByte()))
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			//TODO: handle dispatcher unexpected disconnect
			glg.Fatal(err)
		}
		m := DeserializePeerMessage(message)
		switch m.GetMessage() {
		case CONNFULL:
			w.Disconnect()
			w.Connect()
			break
		case HELLO:
			if w.GetDispatcher() != hex.EncodeToString(m.GetPayload()) {
				w.Disconnect()
				w.Connect()
			}
			w.SetState(INIT)
			glg.Info("P2P: connected to dispatcher")
			break
		case JOB:
			glg.Info("P2P: job received")
			if w.GetState() != LIVE {
				w.SetState(LIVE)
			}
			w.SetBusy(true)
			if m.VerifySignature(w.GetDispatcher()) {
				w.SetItem(qItem.DeserializeItem(m.GetPayload()))
				exec := w.item.Job.Execute(w.item.GetExec(), w.GetDispatcher())
				w.item.SetExec(exec)
				w.conn.WriteMessage(websocket.BinaryMessage, ResultMessage(w.item.GetExec().Serialize(), w.GetPrivByte()))
			} else {
				w.conn.WriteMessage(websocket.BinaryMessage, InvalidSignature())
				w.Disconnect()
			}
			w.SetBusy(false)
			break
		case CANCEL:
			glg.Info("P2P: job cancelled")
			if m.VerifySignature(w.GetDispatcher()) {
				w.item.GetExec().Cancel()
			} else {
				w.conn.WriteMessage(websocket.BinaryMessage, InvalidSignature())
			}
			break
		case SHUT:
			//TODO: handle dispatcher shut
			break
		case SHUTACK:
			for {
				if w.GetBusy() {
					continue
				} else {
					break
				}
			} // wait until worker not busy
			w.Disconnect()
			w.SetState(DOWN)
			glg.Info("Worker: graceful shutdown")
			os.Exit(0)
		default:
			w.Disconnect() //look for new dispatcher
			break
		}
	}
}

func (w Worker) Disconnect() {
	w.conn.Close()
}

func (w *Worker) Connect() {
	for i, dispatcher := range w.GetShortlist() {
		addr, err := ParseAddr(dispatcher)
		if err == nil {
			url := fmt.Sprintf("ws://%v:%v/w", addr["ip"], addr["port"])
			if err = w.Dial(url); err == nil {
				w.SetDispatcher(addr["pub"].(string))
				return
			}
		}
		w.SetShortlist(append(w.GetShortlist()[:i], w.GetShortlist()[i+1:]...))
	}
	w.GetDispatchers()
}

func (w *Worker) Dial(url string) error {
	dailer := websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		ReadBufferSize:  10000,
		WriteBufferSize: 10000,
	}
	conn, _, err := dailer.Dial(url, nil)
	if err != nil {
		glg.Fatal(err)
	}
	conn.EnableWriteCompression(true)
	w.conn = conn
	return nil
}

func (w Worker) WatchInterrupt() {
	select {
	case i := <-w.interrupt:
		glg.Warn("Worker: interrupt detected")
		switch i {
		case syscall.SIGINT, syscall.SIGTERM:
			w.conn.WriteMessage(websocket.BinaryMessage, ShutMessage(w.GetPrivByte()))
			break
		case syscall.SIGQUIT:
			os.Exit(1)
		}
	}
}

func (w *Worker) GetDispatchers() {
	c := NewCentrum()
	res := c.GetDispatchers()
	shortlist, ok := res["dispatchers"]
	if !ok {
		glg.Warn(ErrNoDispatchers)
		os.Exit(0)
	}
	w.SetShortlist(shortlist.([]string))
}

func NewWorker() *Worker {
	core.InitializeDataPath()
	var priv, pub []byte
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// var dbFile string
	// if os.Getenv("ENV") == "dev" {
	// 	dbFile = path.Join(core.IndexPathDev, NodeDB)
	// } else {
	// 	dbFile = path.Join(core.IndexPathProd, NodeDB)
	// }

	// if helpers.FileExists(dbFile) {
	// 	glg.Warn("Worker: using existing keypair")
	// 	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
	// 	if err != nil {
	// 		glg.Fatal(err)
	// 	}
	// 	err = db.View(func(tx *bolt.Tx) error {
	// 		b := tx.Bucket([]byte(NodeBucket))
	// 		priv = b.Get([]byte("priv"))
	// 		pub = b.Get([]byte("pub"))
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		glg.Fatal(err)
	// 	}
	// return &Worker{
	// 	Pub:       pub,
	// 	priv:      priv,
	// 	uptime:    time.Now().Unix(),
	// 	interrupt: interrupt,
	// 	state:     DOWN,
	// 	shortlist: shortlist.([]string),
	// }
	// }
	priv, pub = crypt.GenKeys()
	// db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
	// if err != nil {
	// 	glg.Fatal(err)
	// }

	// err = db.Update(func(tx *bolt.Tx) error {
	// 	b, err := tx.CreateBucket([]byte(NodeBucket))
	// 	if err != nil {
	// 		glg.Fatal(err)
	// 	}

	// 	if err = b.Put([]byte("priv"), priv); err != nil {
	// 		glg.Fatal(err)
	// 	}

	// 	if err = b.Put([]byte("pub"), pub); err != nil {
	// 		glg.Fatal(err)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	glg.Fatal(err)
	// }
	return &Worker{
		Pub:       pub,
		priv:      priv,
		uptime:    time.Now().Unix(),
		interrupt: interrupt,
		state:     DOWN,
	}
}
