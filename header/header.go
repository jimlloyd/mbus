// header.go

package header
// Header data is prefixed to to application message payloads to form mbus packets.

import (
	"bytes"
	"encoding/binary"
)

const (
	Invalid = iota
	Message		// a message packet
	Request 	// a unicast request
	Response 	// a unicast response
	reserved
)

type MessageType uint16  // One of the constants above

var mbusSignature = [4]byte{'m', 'b', 'u', 's'}

type MbusHeader interface {
	Valid() bool
	MessageType() (MessageType, error)
	Encode() (*bytes.Buffer, error)
}

type CommonHeader struct {
	// all fields are specified as byte arrays, with bytes in network byte order (big endian)
	Signature	[4]byte	// a constant signature, 'mbus'
	MsgType		MessageType
}

type MessageHeader struct {
	CommonHeader
	Sequence	uint64	// The sequence number of the first byte of the message
}

type RequestHeader struct {
	CommonHeader
}

type ResponseHeader struct {
	CommonHeader
}

func MakeMessageHeader(sequence uint64) MessageHeader {
	return MessageHeader{CommonHeader{mbusSignature, Message}, sequence}
}

func MakeRequestHeader() RequestHeader {
	return RequestHeader{CommonHeader{mbusSignature, Request}}
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
	return self.Signature == mbusSignature
}

func (self *MessageHeader) Valid() bool {
	return self.Signature == mbusSignature && self.MsgType == Message
}

func (self *RequestHeader) Valid() bool {
	return self.Signature == mbusSignature && self.MsgType == Request
}

func (self *ResponseHeader) Valid() bool {
	return self.Signature == mbusSignature && self.MsgType == Response
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
	return buf, err
}

// Encode the header into a new bytes.Buffer.
// Returns a Buffer so that application message payload can be appended.
func (self *MessageHeader) Encode() (*bytes.Buffer, error) {
	return encodeImpl(self)
}

func (self *RequestHeader) Encode() (*bytes.Buffer, error) {
	return encodeImpl(self)
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


