package p2p

import (
	"encoding/hex"
	"fmt"
	"net"
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
	IP         net.IP
	Port       uint   // port
	Pub        []byte //public key of the node
	Dispatcher string
	priv       []byte //private key of the node
	uptime     int64  //time since node has been up
	conn       *websocket.Conn
	interrupt  chan os.Signal
	shutdown   chan struct{}
	busy       bool
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

func (w Worker) GetIP() net.IP {
	return w.IP
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

func (w Worker) GetPort() int {
	return int(w.Port)
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
	go w.WatchInterrupt()
	w.conn.WriteMessage(websocket.BinaryMessage, HelloMessage(w.GetPubByte()))
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			glg.Fatal(err) //FIXME: error occurs here after job received and processed
		}
		m := DeserializePeerMessage(message)
		switch m.GetMessage() {
		case HELLO:
			w.SetDispatcher(hex.EncodeToString(m.GetPayload()))
			glg.Info("P2P: connected to dispatcher")
			break
		case JOB:
			glg.Info("P2P: job received")
			w.SetBusy(true)
			if m.VerifySignature(w.GetDispatcher()) {
				j := qItem.DeserializeItem(m.GetPayload())
				exec := j.Job.Execute(j.GetExec())
				j.SetExec(exec)
				fmt.Println(string(j.GetExec().Serialize()))
				w.conn.WriteMessage(websocket.BinaryMessage, ResultMessage(j.GetExec().Serialize(), w.GetPrivByte()))
			} else {
				w.conn.WriteMessage(websocket.BinaryMessage, InvalidSignature())
				w.Disconnect()
			}
			w.SetBusy(false)
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

func (w *Worker) Connect(url string) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
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

func NewWorker(port int) *Worker {
	core.InitializeDataPath()
	var priv, pub []byte
	// ip, err := externalip.DefaultConsensus(nil, nil).ExternalIP()
	// if err != nil {
	// 	glg.Fatal(err)
	// }
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
	// 	//FIXME: remove static ws
	// 	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:9999/ws", nil)
	// 	if err != nil {
	// 		glg.Fatal(err)
	// 	}
	// conn.EnableWriteCompression(true)
	// 	return &Worker{
	// 		IP:     ip,
	// 		Pub:    pub,
	// 		priv:   priv,
	// 		Port:   uint(port),
	// 		uptime: time.Now().Unix(),
	// 		conn:   conn,
	// 	}
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
	dailer := websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		ReadBufferSize:  10000,
		WriteBufferSize: 10000,
	}
	conn, _, err := dailer.Dial("ws://127.0.0.1:9999/w", nil)
	if err != nil {
		glg.Fatal(err)
	}
	conn.EnableWriteCompression(true)
	return &Worker{
		// IP:        ip,
		Pub:       pub,
		priv:      priv,
		Port:      uint(port),
		uptime:    time.Now().Unix(),
		conn:      conn,
		interrupt: interrupt,
	}
}
