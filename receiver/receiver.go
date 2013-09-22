package receiver
// A multicast receiver, i.e. subscriber.
// I'm using the name receiver here because this is mostly an experiment while learning Go.
// However, I do intend to implement a message bus subscriber based on what I learn here.

//--------------------------------------------------------------------------------------------------

import (
	"net"
	"fmt"
	"github.com/jimlloyd/mbus/packet"
	"github.com/jimlloyd/mbus/utils"
)

type Receiver struct {
	messageConn *net.UDPConn	// for receiving messages multicast from senders
	controlConn *net.UDPConn	// for sending commands to senders and receiving their responses
	messages    chan packet.Packet
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

	receiver.messages = make(chan packet.Packet, 10)
	go packet.Listen(receiver.messageConn, receiver.messages)

	return receiver, nil
}

func (receiver *Receiver) Close() error {
	err1 := receiver.messageConn.Close()
	err2 := receiver.controlConn.Close()
	close(receiver.messages)
	if (err1 != nil) { return err1 }
	return err2
}

func (receiver *Receiver) MessagesChannel() <-chan packet.Packet {
	return receiver.messages
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
