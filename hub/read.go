package hub

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type readMessage struct {
	ts int64
	id string
}

func encodeRead(id string, ts int64) string {
	id = base64.URLEncoding.EncodeToString([]byte(id))
	return fmt.Sprintf("READ:%s:%d", id, ts)
}

func decodeRead(value string) Message {
	parts := strings.Split(value, ":")
	isRead := len(parts) == 3 && parts[0] == "READ"

	if !isRead {
		return nil
	}

	ts, err := strconv.ParseInt(parts[2], 0, 64)

	if err != nil {
		ts = 0
	}

	id, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &readMessage{
		id: string(id),
		ts: ts,
	}
}

func (s readMessage) Handle(h Hub) error {
	return <-h.Read(s.id, s.ts)
}
