package main

import (
	"fmt"
	"github.com/jimlloyd/mbus/receiver"
)

// Administratively Scoped IP Multicast
// http://tools.ietf.org/html/rfc2365
// 6.2. The IPv4 Organization Local Scope -- 239.192.0.0/14
//   239.192.0.0/14 is defined to be the IPv4 Organization Local Scope,
//   and is the space from which an organization should allocate sub-
//   ranges when defining scopes for private use.

//--------------------------------------------------------------------------------------------------


func main() {

	receiver, err := receiver.NewReceiver("239.192.0.0:5000")
	if err != nil {
		fmt.Println("Error creating receiver:", err)
		return
	}

	incoming := receiver.MessagesChannel()

	for {
		packet := <-incoming
		fmt.Println("Read", len(packet.Data), "bytes:", string(packet.Data), "Remote:",
			packet.Remote())
		
		err := receiver.SendCommand([]byte("Dummy command"), packet.Remote())
		if err != nil {
			fmt.Println("Command not sent, error:", err)
		}
	}
}
