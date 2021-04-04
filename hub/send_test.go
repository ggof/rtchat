package hub

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"
)

const send = "SEND:Z29mZzIzMDE=:cmFuZG9tIG1lc3NhZ2U="

var badSend = []string{
	"SEND1234:Z29mZzIzMDE=:cmFuZG9tIG1lc3NhZ2U=",
	"SEND:Z29mZzIzMDEascasc=:cmFuZG9tIG1lc3NhZ2U=",
	"SEND:Z29mZzIzMDE=:cmFuZG9tIG1lc3NhZ2Uascasc=",
}

func TestDecodeSend(t *testing.T) {
	msg := decodeSend(send)

	if msg == nil {
		t.Error("msg is nil")
	}
}

func TestDecodeSendFail(t *testing.T) {
	for _, s := range badSend {
		msg := decodeSend(s)

		if msg != nil {
			t.Error("msg should be nil")
		}
	}
}

func TestEncodeSend(t *testing.T) {
	str := encodeSend(time.Now().Unix(), "gofg2301", "random message")

	parts := strings.Split(str, ":")

	id, _ := base64.URLEncoding.DecodeString(parts[1])
	msg, _ := base64.URLEncoding.DecodeString(parts[2])

	if string(id) != "gofg2301" {
		t.Error("id does not match")
	}

	if string(msg) != "random message" {
		t.Error("msg does not match")
	}
}

