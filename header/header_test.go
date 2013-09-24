// header_test.go

package header

import (
	"testing"
)

func TestMakeMessageHeader(t *testing.T) {
	const seq = uint32(23)

	h := MakeMessageHeader(seq)

	if !h.Valid() {
		t.Fail()
	}

	msgType, err := h.MessageType()
	if err != nil {
		t.Fail()
	}
	if msgType != Message {
		t.Fail()
	}

	xseq, err := h.GetSequence()
	if err != nil {
		t.Fail()
	}
	if xseq != seq {
		t.Fail()
	}

	buf, err := h.Encode()
	if err != nil {
		t.Fail()
	}

	var x MessageHeader
	buf, err = x.Decode(buf.Bytes())
	if err != nil {
		t.Fail()
	}

	if h != x {
		t.Fail()
	}
}
