package hub 

type Message interface {
	Handle(Hub) error
}

type Decoder = func(string) Message
type Encoder = func(Message) string

var (
	DefaultDecoders = []Decoder{decodeSend, decodeRead}
)
