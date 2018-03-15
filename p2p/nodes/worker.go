package nodes

import (
	"encoding/hex"
	"net"
	"time"
)

type Worker struct {
	IP         net.IP
	Port       uint   // port
	Pub        []byte //public key of the node
	Dispatcher string
	priv       []byte //private key of the node
	uptime     int64  //time since node has been up
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

func (w Worker) GetUptimeString() string {
	return time.Unix(w.uptime, 0).Sub(time.Now()).String()
}

// func NewWorker(port int) *Worker {
// 	core.InitializeDataPath()
// 	if reflect.ValueOf(port).IsNil() {
// 		port = DefaultPort
// 	}
// 	var priv, pub []byte
// 	ip, err := externalip.DefaultConsensus(nil, nil).ExternalIP()
// 	if err != nil {
// 		glg.Fatal(err)
// 	}

// 	var dbFile string
// 	if os.Getenv("ENV") == "dev" {
// 		dbFile = path.Join(core.IndexPathDev, NodeDB)
// 	} else {
// 		dbFile = path.Join(core.IndexPathProd, NodeDB)
// 	}

// 	if helpers.FileExists(dbFile) {
// 		glg.Warn("Dispatcher: using existing keypair and benchmark")
// 		db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
// 		if err != nil {
// 			glg.Fatal(err)
// 		}
// 		err = db.View(func(tx *bolt.Tx) error {
// 			b := tx.Bucket([]byte(NodeBucket))
// 			priv = b.Get([]byte("priv"))
// 			pub = b.Get([]byte("pub"))
// 			return nil
// 		})
// 		if err != nil {
// 			glg.Fatal(err)
// 		}
// 		return &Worker{
// 			IP:     ip,
// 			Pub:    pub,
// 			priv:   priv,
// 			Port:   uint(port),
// 			uptime: time.Now().Unix(),
// 		}
// 	}
// }
