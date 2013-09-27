// sender_receiver_test.go

package receiver

import (
	"testing"
	"time"
	"github.com/jimlloyd/mbus/sender"
)

func MakeReceiver() *Receiver {
	aReceiver, err := NewReceiver("239.192.0.0:5000")
	if err != nil {
		panic("Error creating receiver:" + err.Error())
	}
	return aReceiver
}

func MakeSender() *sender.Sender {
	aSender, err := sender.NewSender("239.192.0.0:5000")
	if err != nil {
		panic("Error creating sender:" + err.Error())
	}
	return aSender
}

func RunReceiver(t *testing.T, aReceiver *Receiver, sem chan<- int, messages []string, numSenders int) {

	receivedMessages := make(map[string]int)
	for _, msg := range(messages) {
		receivedMessages[msg] = 0
	}

	incoming := aReceiver.MessagesChannel()
	for i:=0; i<len(messages)*numSenders; i++ {
		packet := <-incoming
		msg := string(packet.Data)
		count, ok := receivedMessages[msg]
		if !ok {
			t.Error("Unexpected message received:", msg)
		}
		receivedMessages[msg] = count+1
	}

	for msg, count := range(receivedMessages) {
		if count != numSenders {
			t.Error("Wrong number of messages received for message:%s. Expected:%d, received:%d",
				 msg, numSenders, count)
		}
	}

	sem <- 1
}

func RunSender(t *testing.T, aSender *sender.Sender, sem chan<- int, messages []string) {
	for _, expected := range(messages) {
		time.Sleep(50 * time.Millisecond)
		_, err := aSender.Send([]byte(expected))
		if err != nil {
			t.Error("Error sending message:", err)
		}
	}
	sem <- 1
}

func TestSendReceiveNominal(t *testing.T) {
	messages := []string{"aaa", "bbb", "ccccc"}

	const numReceivers = 2
	const numSenders = 3

	senders := []*sender.Sender{}
	receivers := []*Receiver{}

	for i:=0; i<numReceivers; i++ {
		receivers = append(receivers, MakeReceiver())
	}

	for i:=0; i<numSenders; i++ {
		senders = append(senders, MakeSender())
	}

	receiverSem := make(chan int)
	senderSem := make(chan int)

	for _, aReceiver := range(receivers) {
		go RunReceiver(t, aReceiver, receiverSem, messages, numSenders)
	}

	for _, aSender := range(senders) {
		go RunSender(t, aSender, senderSem, messages)
	}

	for i:=0; i<numReceivers; i++ {
		<- receiverSem
	}

	for i:=0; i<numSenders; i++ {
		<- senderSem
	}

	// Can't close here without panic. Seems to only happen here in unit test, not in regular app.
	// I could understand this being a problem if there were multiple receivers, since each listens
	// on the same port, and perhaps the first one to close its connection closes all connections.
	// But this bug happens even when numReceivers=1
	// TODO: figure out the real problem
	// for _, aReceiver := range(receivers) {
	// 	aReceiver.Close()
	// }

	// for _, aSender := range(senders) {
	// 	aSender.Close()
	// }
}
