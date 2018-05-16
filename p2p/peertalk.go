package p2p

const (
	HELLO            = "HELLO"
	INVALIDMESSAGE   = "INVALIDMESSAGE" // invalid message
	CONNFULL         = "CONNFULL"       // max workers reached
	JOB              = "JOB"
	INVALIDSIGNATURE = "JOB"
	RESULT           = "RESULT"
	SHUT             = "SHUT"
	SHUTACK          = "SHUTACK"
	BLOCK            = "BLOCK"
	BLOCKREQ         = "BLOCKREQ"
	BLOCKRES         = "BLOCKRES"
	PEERCONNECT      = "PEERCONNECT"
	PEERDISCONNECT   = "PEERDISCONNECT"
	CANCEL           = "CANCEL"
)

func HelloMessage(payload []byte) []byte {
	return NewPeerMessage(HELLO, payload, nil).Serialize()
}

func InvalidMessage() []byte {
	return NewPeerMessage(INVALIDMESSAGE, nil, nil).Serialize()
}

func ConnFullMessage() []byte {
	return NewPeerMessage(CONNFULL, nil, nil).Serialize()
}

func JobMessage(payload, priv []byte) []byte {
	return NewPeerMessage(JOB, payload, priv).Serialize()
}

func InvalidSignature() []byte {
	return NewPeerMessage(INVALIDMESSAGE, nil, nil).Serialize()
}

func ResultMessage(payload, priv []byte) []byte {
	return NewPeerMessage(RESULT, payload, priv).Serialize()
}

func ShutMessage(priv []byte) []byte {
	return NewPeerMessage(SHUT, nil, priv).Serialize()
}

func ShutAckMessage(priv []byte) []byte {
	return NewPeerMessage(SHUTACK, nil, priv).Serialize()
}

func BlockMessage(payload, priv []byte) []byte {
	return NewPeerMessage(BLOCK, payload, priv).Serialize()
}

func BlockReqMessage(payload, priv []byte) []byte {
	return NewPeerMessage(BLOCKREQ, payload, priv).Serialize()
}

func BlockResMessage(payload, priv []byte) []byte {
	return NewPeerMessage(BLOCKRES, payload, priv).Serialize()
}

func PeerConnectMessage(payload, priv []byte) []byte {
	return NewPeerMessage(PEERCONNECT, payload, priv).Serialize()
}

func PeerDisconnectMessage(payload, priv []byte) []byte {
	return NewPeerMessage(PEERDISCONNECT, payload, priv).Serialize()
}

func CancelMessage(priv []byte) []byte {
	return NewPeerMessage(CANCEL, nil, priv).Serialize()
}
