package p2p

const HELLO = "HELLO"
const INVALIDMESSAGE = "INVALIDMESSAGE" // invalid message
const CONNFULL = "CONNFULL"             // max workers reached

func HelloMessage(payload []byte) PeerMessage {
	return NewPeerMessage(HELLO, payload, nil)
}

func InvalidMessage() PeerMessage {
	return NewPeerMessage(INVALIDMESSAGE, nil, nil)
}

func ConnFull() PeerMessage {
	return NewPeerMessage(CONNFULL, nil, nil)
}

// func AckMessage() PeerMessage {
// 	return PeerMessage{Message: "Ack"}
// }
