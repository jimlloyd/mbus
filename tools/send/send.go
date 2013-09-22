package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"
	"github.com/jimlloyd/mbus/sender"
)

//--------------------------------------------------------------------------------------------------

// Administratively Scoped IP Multicast
// http://tools.ietf.org/html/rfc2365
// 6.2. The IPv4 Organization Local Scope -- 239.192.0.0/14
//   239.192.0.0/14 is defined to be the IPv4 Organization Local Scope,
//   and is the space from which an organization should allocate sub-
//   ranges when defining scopes for private use.

func main() {

	var message []byte;
	var err error;

	flag.Parse()

	if flag.NArg() == 0 {
		const layout = "Jan 2, 2006 at 3:04:05pm (MST)"
		message = []byte(time.Now().Format(layout))
		fmt.Println("Using default message:", string(message))
	} else {
		fmt.Println("Reading contents of file:", flag.Arg(0))
		message, err = ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	}

	sender, err := sender.NewSender("239.192.0.0:5000")
	if err != nil {
		fmt.Println("Error creating sender:", err)
		return
	}

	defer sender.Close()

	for i := 0; i < 10; i++ {
		time.Sleep(1000 * time.Millisecond)
		nbytes, err := sender.Send(message)
		if err != nil {
			fmt.Println("Error sending:", err)
		}
		fmt.Println("Wrote", nbytes, "bytes.")
	}
}
