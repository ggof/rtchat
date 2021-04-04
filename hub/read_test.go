package hub

import (
	"testing"
)

const read = "READ:Z29mZzIzMDE=:1617483677"

func TestDecodeRead(t *testing.T) {

	msg := decodeRead(read)

	if msg == nil {
		t.Error("msg is nil")
	}
}
