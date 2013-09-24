package packet

//--------------------------------------------------------------------------------------------------

import (
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

