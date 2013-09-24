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

var mbusSignature = [4]byte{'m', 'b', 'u', 's'}

type CommonHeader struct {
	// all fields are specified as byte arrays, with bytes in network byte order (big endian)
	Signature	[4]byte	// a constant signature, 'mbus'
	MsgType		uint16  // One of the constants above
}

type MessageHeader struct {
	CommonHeader
	Sequence	uint32	// The sequence number of the first byte of the message
}

func MakeMessageHeader(sequence uint32) MessageHeader {
	return MessageHeader{CommonHeader{mbusSignature, Message}, sequence}
}

func (self *CommonHeader) Valid() bool {
	return self.Signature == mbusSignature
}

func (self *MessageHeader) Valid() bool {
	return self.Signature == mbusSignature && self.MsgType == Message
}

func (self *CommonHeader) MessageType() (uint16, error) {
	if !self.Valid() {
		return Invalid, InvalidHeaderError{}
	}
	return self.MsgType, nil
}

func (self *MessageHeader) GetSequence() (uint32, error) {
	if !self.Valid() || self.MsgType!=Message {
		return Invalid, InvalidHeaderError{}
	}
	return self.Sequence, nil
}

// Encode the header into a new bytes.Buffer.
// Returns a Buffer so that application message payload can be appended.
func (self *MessageHeader) Encode() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if !self.Valid() || self.MsgType!=Message {
		return buf, InvalidHeaderError{}
	}
	err := binary.Write(buf, binary.LittleEndian, *self)
	return buf, err
}

// Decode the header from the bytes slice into this MessageHeader.
// Returns the Buffer so that application message payload can be retrieved.
func (self *MessageHeader) Decode(packetData []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(packetData)
	err := binary.Read(buf, binary.LittleEndian, self)
	if err==nil {
		if !self.Valid() || self.MsgType!=Message {
			return buf, InvalidHeaderError{}
		}
	}
	return buf, err
}

type InvalidHeaderError struct {
}

func (InvalidHeaderError) Error() string {
	return "Not a valid mbus header"
}


