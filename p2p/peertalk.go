package p2p

const (
	HELLO               = "HELLO"
	INVALIDMESSAGE      = "INVALIDMESSAGE" // invalid message
	CONNFULL            = "CONNFULL"       // max workers reached
	JOB                 = "JOB"
	INVALIDSIGNATURE    = "JOB"
	RESULT              = "RESULT"
	SHUT                = "SHUT"
	SHUTACK             = "SHUTACK"
	VERSION             = "VERSION"
	BLOCK               = "BLOCK"
	NEIGHBOURS          = "NEIGHBOURS"
	NEIGHBOURCONNECT    = "NEIGHBOURCONNECT"
	NEIGHBOURDISCONNECT = "NEIGHBOURDISCONNECT"
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

func VersionMessage(payload, priv []byte) []byte {
	return NewPeerMessage(VERSION, payload, priv).Serialize()
}

func BlockMessage(payload, priv []byte) []byte {
	return NewPeerMessage(BLOCK, payload, priv).Serialize()
}

func NeighbourConnectMessage(payload, priv []byte) []byte {
	return NewPeerMessage(NEIGHBOURCONNECT, payload, priv).Serialize()
}

func NeighbourDisconnectMessage(payload, priv []byte) []byte {
	return NewPeerMessage(NEIGHBOURDISCONNECT, payload, priv).Serialize()
}

func NeighboursMessage(payload []byte, priv []byte) []byte {
	return NewPeerMessage(NEIGHBOURS, payload, priv).Serialize()
}
