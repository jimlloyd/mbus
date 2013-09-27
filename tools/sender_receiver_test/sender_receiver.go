package main

import (
	"fmt"
	"time"
	"github.com/jimlloyd/mbus/receiver"
	"github.com/jimlloyd/mbus/sender"
)

func MakeReceiver() *receiver.Receiver {
	aReceiver, err := receiver.NewReceiver("239.192.0.0:5000")
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

func RunReceiver(aReceiver *receiver.Receiver, sem chan<- int, messages []string, numSenders int) {

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
			fmt.Println("Unexpected message received:", msg)
		}
		receivedMessages[msg] = count+1
	}

	for msg, count := range(receivedMessages) {
		if count != numSenders {
			fmt.Printf("Wrong number of messages received for message:%s. Expected:%d, received:%d",
				 msg, numSenders, count)
		}
	}

	fmt.Println("Done receiving")
	sem <- 1
}

func RunSender(aSender *sender.Sender, sem chan<- int, messages []string) {
	for _, expected := range(messages) {
		fmt.Println("Sending message:", expected)
		time.Sleep(50git  * time.Millisecond)
		_, err := aSender.Send([]byte(expected))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
		fmt.Println("Sent message:", expected)
	}
	fmt.Println("Done sending")
	sem <- 1
}

func main() {
	messages := []string{"aaa", "bbb", "ccccc"}

	const numReceivers = 3
	const numSenders = 5

	senders := []*sender.Sender{}
	receivers := []*receiver.Receiver{}

	for i:=0; i<numReceivers; i++ {
		receivers = append(receivers, MakeReceiver())
	}

	for i:=0; i<numSenders; i++ {
		senders = append(senders, MakeSender())
	}
	receiverSem := make(chan int)
	senderSem := make(chan int)

	for _, aReceiver := range(receivers) {
		go RunReceiver(aReceiver, receiverSem, messages, numSenders)
	}

	for _, aSender := range(senders) {
		go RunSender(aSender, senderSem, messages)
	}

	for i:=0; i<numReceivers; i++ {
		<- receiverSem
	}

	for i:=0; i<numSenders; i++ {
		<- senderSem
	}

	for _, aReceiver := range(receivers) {
		aReceiver.Close()
	}

	for _, aSender := range(senders) {
		aSender.Close()
	}
}
