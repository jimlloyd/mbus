// packet_test.go

package packet

import (
	"fmt"
	"testing"
)

func droperTestAtRate(t *testing.T, droprate float64) {
	c := make(chan []byte)

	total := 1000	// push 1000 messages total

	sem := make(chan int)
	count := 0
	go func(in <-chan []byte) {
		for {
			m := <- in
			if m == nil {
				break
			}
			count++
		}
		sem <- 1
		return
	} (Dropper(c, droprate))

	for i:=0; i<total; i++ {
		msg := []byte(fmt.Sprintf("Msg%d", i))
		c <- msg
	}

	c <- nil
	<- sem

	actualRate := float64(total - count) / float64(total)
	fmt.Println("Expected rate:", droprate, "Actual rate:", actualRate)
	if actualRate < droprate-0.02 {
		t.Error("Drop rate appears too low")
	}
	if actualRate > droprate+0.02 {
		t.Error("Drop rate appears too high")
	}
}

func TestDropper(t *testing.T) {
	droperTestAtRate(t, 0.2)
	droperTestAtRate(t, 0.5)
	droperTestAtRate(t, 0.8)
}
