package p2p

const HELLO = "HELLO"
const INVALIDMESSAGE = "INVALIDMESSAGE" // invalid message
const CONNFULL = "CONNFULL"             // max workers reached
const JOB = "JOB"
const INVALIDSIGNATURE = "JOB"
const RESULT = "RESULT"

func HelloMessage(payload []byte) []byte {
	return NewPeerMessage(HELLO, payload, nil).Serialize()
}

func InvalidMessage() []byte {
	return NewPeerMessage(INVALIDMESSAGE, nil, nil).Serialize()
}

func ConnFullMessage() []byte {
	return NewPeerMessage(CONNFULL, nil, nil).Serialize()
}

func JobMessage(payload []byte, priv []byte) []byte {
	return NewPeerMessage(JOB, payload, priv).Serialize()
}

func InvalidSignature() []byte {
	return NewPeerMessage(INVALIDMESSAGE, nil, nil).Serialize()
}

func ResultMessage(payload []byte, priv []byte) []byte {
	return NewPeerMessage(RESULT, payload, priv).Serialize()
}
