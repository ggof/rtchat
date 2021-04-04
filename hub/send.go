package hub

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

type sendMessage struct {
	ts  int64
	id  string
	msg string
}

func decodeSend(value string) Message {
	parts := strings.Split(value, ":")
	isSend := len(parts) == 3 && parts[0] == "SEND"

	if !isSend {
		return nil
	}

	id, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	msg, err := base64.URLEncoding.DecodeString(parts[2])
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &sendMessage{
		id:  string(id),
		msg: string(msg),
	}
}

func encodeSend(ts int64, id, msg string) string {
	id = base64.URLEncoding.EncodeToString([]byte(id))
	msg = base64.URLEncoding.EncodeToString([]byte(msg))
	return fmt.Sprintf("SEND:%s:%s:%d", id, msg, ts)
}

func (s sendMessage) Handle(h Hub) error {
	log.Printf("%s sent message %s\n", s.id, s.msg)
	return <-h.Send(s.id, s.msg)
}
