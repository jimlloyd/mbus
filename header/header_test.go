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

	x.Signature = [4]byte{23} // just some invalid signature
	if x.Valid() {
		t.Error("Header with invalid Signature is considered valid")
	}

	_, err = x.Encode()
	if err == nil {
		t.Error("Header with invalid Signature encodes without returning an error")
	}

	x.Signature = mbusSignature // restore the signature so following tests only test MsgType

	x.MsgType = 99 // just some invalid value
	if x.Valid() {
		t.Error("Header with invalid MsgType is considered valid")
	}

	_, err = x.Encode()
	if err == nil {
		t.Error("Header with invalid MsgType encodes without returning an error")
	}
}
