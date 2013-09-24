// header.go

package header
// Header is the data prefixed to all mbus packets.

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

func (self *CommonHeader) MessageType() uint16 {
	return self.MsgType
}

func (self *MessageHeader) GetSequence() uint32 {
	return self.Sequence
}

