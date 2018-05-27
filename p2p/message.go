package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/kpango/glg"
)

//PeerMessage message sent between nodes
type PeerMessage struct {
	Message   string
	Payload   []byte
	Signature [][]byte
}

func NewPeerMessage(message string, payload []byte, priv []byte) PeerMessage {
	pm := PeerMessage{Message: message, Payload: payload}
	if priv != nil {
		pm.sign(priv)
	}
	return pm
}

func (m PeerMessage) GetMessage() string {
	return m.Message
}

func (m PeerMessage) GetPayload() []byte {
	return m.Payload
}

func (m PeerMessage) GetSignature() [][]byte {
	return m.Signature
}

func (m *PeerMessage) SetMessage(message string) {
	m.Message = message
}

func (m *PeerMessage) SetPayload(payload []byte) {
	m.Payload = payload
}

func (m *PeerMessage) setSignature(sig [][]byte) {
	m.Signature = sig
}

func (m *PeerMessage) sign(priv []byte) {
	hash := sha256.Sum256(bytes.Join(
		[][]byte{
			[]byte(m.GetMessage()),
			m.GetPayload(),
		},
		[]byte{},
	))
	privateKey, _ := x509.ParseECPrivateKey(priv)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		glg.Fatal("Unable to sign peer message")
	}
	var temp [][]byte
	temp = append(temp, r.Bytes(), s.Bytes())
	m.setSignature(temp)
}

func (m *PeerMessage) VerifySignature(pub string) bool {
	pubBytes, err := hex.DecodeString(pub)
	if err != nil {
		glg.Fatal(err)
	}
	var r big.Int
	var s big.Int
	r.SetBytes(m.GetSignature()[0])
	s.SetBytes(m.GetSignature()[1])

	publicKey, _ := x509.ParsePKIXPublicKey(pubBytes)
	hash := sha256.Sum256(bytes.Join(
		[][]byte{
			[]byte(m.GetMessage()),
			m.GetPayload(),
		},
		[]byte{},
	))
	switch pubConv := publicKey.(type) {
	case *ecdsa.PublicKey:
		return ecdsa.Verify(pubConv, hash[:], &r, &s)
	default:
		return false
	}
}

func (m PeerMessage) Serialize() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		glg.Fatal(err)
	}
	return bytes
}

func DeserializePeerMessage(b []byte) PeerMessage {
	var temp PeerMessage
	err := json.Unmarshal(b, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}
