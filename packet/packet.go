package packet

//--------------------------------------------------------------------------------------------------

import (
	"math/rand"
	"net"
)

type Packet struct {
	Data   []byte
	remote net.Addr
}

func (packet Packet) Remote() net.Addr {
	return packet.remote
}

func Listen(conn *net.UDPConn, incoming chan<- Packet) error {
	for {
		data := make([]byte, 8192)
		size, remote, err := conn.ReadFrom(data)
		if err != nil {
			panic(err)
		}
		incoming <- Packet{data[0:size], remote}
	}
}

func Dropper(payloads <-chan []byte, droprate float64) <-chan []byte {
	// Note, as a special case never drop nil messages. Currently used only in testing.
	filtered := make(chan []byte)
	go func() {
		for {
			msg := <- payloads
			if msg==nil || rand.Float64() >= droprate {
				filtered <- msg
			}
		}
	} ()
	return filtered
}

