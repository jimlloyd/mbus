package sender
// A multicast sender, i.e. publisher.
// I'm using the name sender here because this is mostly an experiment while learning Go.
// However, I do intend to implement a message bus publisher based on what I learn here.

//--------------------------------------------------------------------------------------------------

import (
	"net"
	"fmt"
	"github.com/jimlloyd/mbus/packet"
	"github.com/jimlloyd/mbus/utils"
)

type Sender struct {
	conn *net.UDPConn
	mcast *net.UDPAddr
}

func NewSender(mcastAddress string) (*Sender, error) {
	var err error

	sender := new(Sender)

	// One connection is used both to send multicasts and to receive command packets.
	// We create the connection by setting up listening for commands packets,
	// but can also use the connection to send multicasts.
	sender.conn, err = utils.ListenUDP4()
	if err != nil {
		return nil, err
	}

	sender.mcast, err = net.ResolveUDPAddr("udp4", mcastAddress)
	if err != nil {
		return nil, err
	}

	commands := make(chan packet.Packet, 10)
	go packet.Listen(sender.conn, commands)
	go sender.serveCommand(commands)

	return sender, nil
}

func (sender *Sender) Close() error {
	return sender.conn.Close()
}

func (sender *Sender) Send(message []byte) (int, error) {
	return sender.conn.WriteToUDP(message, sender.mcast)
}

func (sender *Sender) serveCommand(commands <-chan packet.Packet) {
	for {
		packet := <-commands

		// We'll eventually respond to meaningful commands, but for now just log them.
		fmt.Println("Received command", len(packet.Data()), "bytes:", string(packet.Data()),
			"Remote:", packet.Remote())
	}
}


