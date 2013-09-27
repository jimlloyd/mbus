// sendersmap.go
// A map of the senders that a receiver is receiving messages from.

package sendersmap

import (
	"sync"
	"github.com/jimlloyd/mbus/packet"
)
type SenderInfo struct {
	Addr 	string

	// a count of packets received
	// we don't really care about the count, but it's useful now for development/debugging
	Count	int

	// The sequence number we next expect to deliver.
	// When we receive that packet, we can deliver it immediately,
	// and update this value by the packet payload length
	DeliveredTo uint64

	// ----- Notes about out of order packet handling.
	// -- If we receive a packet whose sequence number is less than DeliverdTo
	// there are two possibilities:
	// 1. We are seeing a packet duplicate. In this case we drop the packet.
	// 2. The process that was sending was restarted and the exact same sender port
	//    was reused. In this case, we want to keep the packet, and must reset
	//    ReceivedTo and DeliveredTo.
	// Both of these scenarios are unlikely but possible and must be handled.
	// The first scenario is probably several orders of magnitude more likely.
	// -- If we receive a packet whose sequence number is greater than DeliveredTo
	// then we have apparently missed one or more packets. Again there are two possibilites:
	// 1. DeliveredTo is zero, which means we haven't seen any previous packets from this
	//    sender. In that case we deliver the packet and update DeliveredTo to the sequence
	//    number following this packet. That leaves the possibility that we may receive
	//    resent/duplicate packets with lower sequence numbers, which we will drop even
	//    though we haven't delivered them.
	// 2. DeliveredTo is nonzero. This indicates the expected packet was dropped
	//    or delayed. We need to hold this packet for later delivery, and may need to
	//    notify sender to resend the missing range of bytes. It is best to ask for the
	//    range of bytes, since it is possible that multiple packets were dropped or delayed.

	Holding map[uint64]packet.Packet

	// other fields we'll probably track:
	// 1. Last message tick
}

type SendersMap struct {
	rep		map[string]*SenderInfo
	lock	sync.RWMutex
}

func New() *SendersMap {
	sendersMap := new(SendersMap)
	sendersMap.rep = map[string]*SenderInfo{}
	return sendersMap
}

func (self *SendersMap) Get(addr string) *SenderInfo {
	self.lock.RLock()
	info, ok := self.rep[addr]
	self.lock.RUnlock()

	if ok {
		return info
	}

	self.lock.Lock()
	info, ok = self.rep[addr]
	if !ok {
		info = &SenderInfo{addr, 0, 0, make(map[uint64]packet.Packet)}
		self.rep[addr] = info
	}
	self.lock.Unlock()
	return info
}
