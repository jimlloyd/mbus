// header_test.go

package header

import (
	"testing"
)

func TestHeaderTypes(t *testing.T) {
	const seq = uint64(23)
	message := MakeMessageHeader(seq)
	request := MakeRequestHeader(MakeFixedSignature("Resend.."))
	response := MakeResponseHeader()

	if !message.Valid() {
		t.Fail()
	}
	if !request.Valid() {
		t.Fail()
	}
	if !response.Valid() {
		t.Fail()
	}

	messageBuf, err := message.Encode()
	if err != nil {
		t.Fail()
	}
	requestBuf, err := request.Encode()
	if err != nil {
		t.Fail()
	}
	responseBuf, err := response.Encode()
	if err != nil {
		t.Fail()
	}

	var xMessage MessageHeader
	var xRequest RequestHeader
	var xResponse ResponseHeader

	_, err = xMessage.Decode(messageBuf.Bytes())
	if err != nil {
		t.Fail()
	}
	if message != xMessage {
		t.Fail()
	}

	_, err = xRequest.Decode(requestBuf.Bytes())
	if err != nil {
		t.Fail()
	}
	if request != xRequest {
		t.Fail()
	}

	_, err = xResponse.Decode(responseBuf.Bytes())
	if err != nil {
		t.Fail()
	}
	if response != xResponse {
		t.Fail()
	}

	if Message != PeekMessageType(messageBuf.Bytes()) {
		t.Error("PeekMessageType failed to return Message for a Message packet.")
	}
	if Request != PeekMessageType(requestBuf.Bytes()) {
		t.Error("PeekMessageType failed to return Request for a Message packet.")
	}
	if Response != PeekMessageType(responseBuf.Bytes()) {
		t.Error("PeekMessageType failed to return Response for a Message packet.")
	}


	_, err = xResponse.Decode(requestBuf.Bytes())
	if err == nil {
		t.Error("A ResponseHeader should not accept decoding from a Request packet")
	}

	_, err = xMessage.Decode(requestBuf.Bytes())
	if err == nil {
		t.Error("A MessageHeader should not accept decoding from a Request packet")
	}
}


func TestMakeMessageHeader(t *testing.T) {
	const seq = uint64(23)
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

	x.MbusSig = Signature{23} // just some invalid signature
	if x.Valid() {
		t.Error("Header with invalid Signature is considered valid")
	}

	_, err = x.Encode()
	if err == nil {
		t.Error("Header with invalid Signature encodes without returning an error")
	}

	x.MbusSig = mbusSignature // restore the signature so following tests only test MsgType

	x.MsgType = 99 // just some invalid value
	if x.Valid() {
		t.Error("Header with invalid MsgType is considered valid")
	}

	_, err = x.Encode()
	if err == nil {
		t.Error("Header with invalid MsgType encodes without returning an error")
	}
}
