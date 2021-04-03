package messages

const (
	connect    = "CONNECT" // new person connects
	disconnect = "DISCONNECT" // someone disconnected
	send       = "SEND" // someone sends a message
	receive    = "RECEIVE" // someone has read a message
)

// Message struct representing user interaction
type Message struct {
	TS     int64
	SentBy string
	ReadBy []string
	Value  string
}

// Protocol goes as follows: 
// CONNECT and DISCONNECT have no Msg
// SEND and RECEIVE have a Msg
type WSMessage struct {
	Type string `json:"type"`
	Msg  *Message `json:"msg,omitempty"`
}

