// receiver.go
package receiver
// A multicast receiver, i.e. subscriber.
// I'm using the name receiver here because this is mostly an experiment while learning Go,
// but also because this is a primitive form of subscriber. The ultimate subscriber
// may be a layer on top of receiver.
//--------------------------------------------------------------------------------------------------

import (
	"bytes"
	"net"
	"fmt"
	"github.com/jimlloyd/mbus/header"
	"github.com/jimlloyd/mbus/packet"
	"github.com/jimlloyd/mbus/utils"
	"github.com/jimlloyd/mbus/receiver/sendersmap"
)

type Receiver struct {
	messageConn *net.UDPConn	// for receiving messages multicast from senders
	controlConn *net.UDPConn	// for sending commands to senders and receiving their responses

	incoming    chan packet.Packet 	// message packets received but not yet analyzed/sequenced
	sequenced	chan packet.Packet  // message packets sequenced and ready for application to process

	senders 	*sendersmap.SendersMap
}

func NewReceiver(mcastAddress string) (*Receiver, error) {
	receiver := new(Receiver)

	addr, err := net.ResolveUDPAddr("udp4", mcastAddress)
	if err != nil {
		return nil, err
	}

	receiver.messageConn, err = net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	receiver.controlConn, err = utils.ListenUDP4()
	if err != nil {
		receiver.messageConn.Close()
		return nil, err
	}
	fmt.Println("Listening for command/control responses on local address:", receiver.controlConn.LocalAddr())

	// Currently there is just one channel that delivers all messages as they are received.
	// TODO:
	// 1. Track unique senders
	// 2. Track sequence numbers of packets from each sender
	// 3. Detect dropped packets and send request for packet resend to correct sender.
	// 4. Detect duplicate packets and drop them.
	// 5. Optionally hold future packets for delivery in correct sequence.

	receiver.senders = sendersmap.New()

	receiver.incoming = make(chan packet.Packet, 10)
	receiver.sequenced = make(chan packet.Packet, 10)

	go receiver.AnalyzeAndSequence()
	go packet.Listen(receiver.messageConn, receiver.incoming)

	return receiver, nil
}

func (receiver *Receiver) Close() error {
	err1 := receiver.messageConn.Close()
	err2 := receiver.controlConn.Close()
	close(receiver.incoming)
	close(receiver.sequenced)
	if (err1 != nil) { return err1 }
	return err2
}

func (receiver *Receiver) MessagesChannel() <-chan packet.Packet {
	return receiver.sequenced
}

type TruncatedError struct {
}

func (TruncatedError) Error() string {
	return "Truncated Packet"
}

func (receiver *Receiver) SendCommand(command []byte, addr net.Addr) error {
	l, err := receiver.controlConn.WriteTo(command, addr)
	if err==nil && l!=len(command) {
		return TruncatedError{}
	}
	return err
}

func (receiver *Receiver) AnalyzeAndSequence() {
	for {
		packet := <-receiver.incoming
		fmt.Println("Analyzing a packet")

		senderInfo := receiver.senders.Get(packet.Remote().String())
		senderInfo.Count++
		fmt.Println("Received", senderInfo.Count, "packets from", senderInfo.Addr)

		var head header.MessageHeader

		buf := bytes.NewBuffer(packet.Data)
		buf, err := head.Decode(buf.Bytes())
		if err != nil {
			fmt.Println("Dropping invalid packet. Header:", head, "Error:", err)
		} else {
			packet.Data = buf.Bytes()
			payloadLen := uint64(len(packet.Data))

			// this is just a hack for now
			// TODO: check for missing packets
			nextSeq := head.Sequence + payloadLen
			senderInfo.ReceivedTo = nextSeq

			receiver.sequenced <- packet
		}

	}
}
