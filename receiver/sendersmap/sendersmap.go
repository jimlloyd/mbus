// sendersmap.go
// A map of the senders that a receiver is receiving messages from.

package sendersmap

import (
	"sync"
	"fmt"
)

type SenderInfo struct {
	Addr 	string

	// a count of packets received
	// we don't really care about the count, but it's useful now for development/debugging
	Count	int

	ReceivedTo	uint64

	// other fields we'll probably track:
	// 1. Sequence number delivered to
	// 2. Sequence number received to
	// 3. Last message tick
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
	fmt.Println("Looking up", addr, "in map with", len(self.rep), "elements")

	self.lock.RLock()
	info, ok := self.rep[addr]
	self.lock.RUnlock()

	if ok {
		fmt.Println("Address", addr, "already in map:", info)
		return info
	}

	self.lock.Lock()
	info, ok = self.rep[addr]
	if !ok {
		info = &SenderInfo{addr, 0, 0}
		self.rep[addr] = info
		fmt.Println("Adding new address", addr, "to map:", info)
	}
	self.lock.Unlock()
	return info
}