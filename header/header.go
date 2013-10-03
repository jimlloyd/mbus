// header.go

package header
// Header data is prefixed to to application message payloads to form mbus packets.

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	Invalid = iota
	Message		// a message packet
	Request 	// a unicast request
	Response 	// a unicast response
	reserved
)

type MessageType uint16  // One of the constants above

const SignatureSize = 8
type Signature [SignatureSize]byte 	 // a 4 byte 'signature' used in message encodings to indicate type/kind of messages

var mbusSignature = MakeFixedSignature("gobusgo!")

type MbusHeader interface {
	Valid() bool
	MessageType() (MessageType, error)
	Encode() (*bytes.Buffer, error)
}

type CommonHeader struct {
	// all fields are specified as byte arrays, with bytes in network byte order (big endian)
	MbusSig		Signature	// a constant signature ('mbus') used to provide confidence that message is valid
	MsgType		MessageType
}

type MessageHeader struct {
	CommonHeader
	Sequence	uint64	// The sequence number of the first byte of the message
}

type RequestHeader struct {
	CommonHeader
	Verb		Signature
}

type ResponseHeader struct {
	CommonHeader
}

func MakeMessageHeader(sequence uint64) MessageHeader {
	return MessageHeader{CommonHeader{mbusSignature, Message}, sequence}
}

func MakeRequestHeader(verb Signature) RequestHeader {
	return RequestHeader{CommonHeader{mbusSignature, Request}, verb}
}

func MakeResponseHeader() ResponseHeader {
	return ResponseHeader{CommonHeader{mbusSignature, Response}}
}

func PeekMessageType(packetData []byte) MessageType {
	var head CommonHeader
	buf := bytes.NewBuffer(packetData)
	err := binary.Read(buf, binary.LittleEndian, &head)
	if err == nil {
		switch head.MsgType {
		case Message, Request, Response:
			return head.MsgType
		}
	}
	return Invalid
}

func (self *CommonHeader) Valid() bool {
	return self.MbusSig == mbusSignature
}

func (self *MessageHeader) Valid() bool {
	return self.MbusSig == mbusSignature && self.MsgType == Message
}

func (self *RequestHeader) Valid() bool {
	return self.MbusSig == mbusSignature && self.MsgType == Request
}

func (self *ResponseHeader) Valid() bool {
	return self.MbusSig == mbusSignature && self.MsgType == Response
}

func (self *CommonHeader) MessageType() (MessageType, error) {
	if !self.Valid() {
		return Invalid, InvalidHeaderError{}
	}
	return self.MsgType, nil
}

func (self *MessageHeader) GetSequence() (uint64, error) {
	if !self.Valid() || self.MsgType!=Message {
		return Invalid, InvalidHeaderError{}
	}
	return self.Sequence, nil
}

func encodeImpl(self MbusHeader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if !self.Valid() {
		return buf, InvalidHeaderError{}
	}
	err := binary.Write(buf, binary.LittleEndian, self)
	if err != nil {
		fmt.Println("encodeImp failed with err:", err)
	}
	return buf, err
}

// Encode the header into a new bytes.Buffer.
// Returns a Buffer so that application message payload can be appended.
func (self *MessageHeader) Encode() (*bytes.Buffer, error) {
	return encodeImpl(self)
}

func (self *RequestHeader) Encode() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if !self.Valid() {
		return buf, InvalidHeaderError{}
	}
	err := binary.Write(buf, binary.LittleEndian, self)
	if err != nil {
		fmt.Println("RequestHeader.Encode() failed with err:", err)
	}
	return buf, err
}

func (self *ResponseHeader) Encode() (*bytes.Buffer, error) {
	return encodeImpl(self)
}

func decodeImpl(self MbusHeader, packetData []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(packetData)
	err := binary.Read(buf, binary.LittleEndian, self)
	if err==nil {
		if !self.Valid() {
			return buf, InvalidHeaderError{}
		}
	}
	return buf, err
}

// Decode the header from the bytes slice into this MessageHeader.
// Returns the Buffer so that application message payload can be retrieved.
func (self *MessageHeader) Decode(packetData []byte) (*bytes.Buffer, error) {
	return decodeImpl(self, packetData)
}

func (self *RequestHeader) Decode(packetData []byte) (*bytes.Buffer, error) {
	return decodeImpl(self, packetData)
}

func (self *ResponseHeader) Decode(packetData []byte) (*bytes.Buffer, error) {
	return decodeImpl(self, packetData)
}

type InvalidHeaderError struct {
}

func (InvalidHeaderError) Error() string {
	return "Not a valid mbus header"
}

func MakeRequest(verb Signature, parameters []byte) ([]byte, error) {
	h := MakeRequestHeader(verb)
	buf, err := h.Encode()
	if err != nil {
		fmt.Println("Failed to encode request:", err)
		return nil, err
	}
	buf.Write(parameters)
	return buf.Bytes(), nil
}

func MakeFixedSignature(s string) Signature {
	var sig Signature
	if len(s) != SignatureSize {
		panic("Convenience function MakeFixedSignature must be called with ascii string exactly 8 chars long.")
	}
	for i:=0; i<SignatureSize; i++ {
		sig[i] = s[i]
	}
	return sig
}

