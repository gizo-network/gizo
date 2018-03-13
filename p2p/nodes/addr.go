package nodes

import (
	"encoding/hex"
	"errors"
	"net"
	"net/url"
)

type nodeInterface interface {
	GetIP() net.IP
	GetRPCPort() int
	GetWSPort() int
	GetPubString() string
	NodeTypeDispatcher() bool
}

var (
	ErrInvalidIP       = errors.New("Unable to parse IP string")
	ErrInvalidNodeType = errors.New("Invalide node type")
)

func NodeAddr(n nodeInterface) string {
	u := url.URL{}
	if n.NodeTypeDispatcher() {
		u.Scheme = DispatcherScheme
	} else {
		u.Scheme = WorkerScheme
	}
	addr := net.TCPAddr{IP: n.GetIP(), Port: n.GetWSPort()}
	u.User = url.User(n.GetPubString())
	u.Host = addr.String()
	return u.String()
}

func ParseDispatcher(raw string) (*Dispatcher, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if u.Scheme != DispatcherScheme {
		return nil, ErrInvalidNodeType
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, ErrInvalidIP
	}

	pub, err := hex.DecodeString(u.User.String())
	if err != nil {
		return nil, err
	}

	return &Dispatcher{
		IP:   net.ParseIP(host),
		Port: port,
		Pub:  pub,
	}, nil
}
