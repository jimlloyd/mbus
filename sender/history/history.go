// mbus/sender/history/history.go

package history

import (
	"time"
)

type SeqNum     uint64	// an absolute position in a range of bytes
type SizeBytes	uint64	// the size of a span of bytes (delta between sequence)

type Tick 		int64	// an absolute tick
type AgeNanos   int64	// a duration of time (delta between ticks)

type History struct {
	// The sorted list of sequence numbers currently in the history
	sequences []SeqNum

	// The messages in our history.
	// The map key is the sequence number (from MessageHeader.Sequence).
	// The map value is the complete message (MessageHeader+payload) ready as previously sent
	messages map[SeqNum][]byte

	// For each message, the tick at the time it was added to this history (essentially time it was sent)
	ticks map[SeqNum]Tick

	// Keep history for at least this duration, even when maxStorage is exceeded
	minAge AgeNanos

	// Purge oldest messages when their age exceeds this maximum nanoseconds
	maxAge	AgeNanos

	// Purge oldest messages when total exceeds this maximum, unless minAge applies
	maxStorage SizeBytes
}

func NewHistory(minAgeSeconds int32, maxAgeSeconds int32, maxPayloadMB uint32) (*History) {
	if minAgeSeconds > maxAgeSeconds {
		panic("History minAge must be less than maxAge")
	}

	minAge := AgeNanos(minAgeSeconds)*10e9
	maxAge := AgeNanos(maxAgeSeconds)*10e9
	maxStorage := SizeBytes(maxPayloadMB)*10e6

	return &History{[]SeqNum{}, make(map[SeqNum][]byte), make(map[SeqNum]Tick), minAge, maxAge, maxStorage }
}

func (self *History) Add(seq uint64, message []byte) {
	sequence := SeqNum(seq)
	tick := Tick(time.Now().UnixNano())

	prevLen := len(self.sequences)
	if (prevLen>0 && self.sequences[prevLen-1]>=sequence) {
		panic("Sequence numbers must be increasing")
	}
	self.sequences = append(self.sequences, sequence)
	self.messages[sequence] = message
	self.ticks[sequence] = tick

	self.purgeOldest()
}

func (self *History) Length() int {
	length := len(self.sequences)
	if length != len(self.messages) { panic("Inconsistent History, different length for messages.") }
	if length != len(self.ticks) { panic("Inconsistent History, different length for ticks.") }
	return length
}

// If there is a message in the history with the given seqNum, return it.
// Otherwise return nil. For now we don't try to give any hints as to why the lookup failed.
// If the caller provides a previously used sequence number then perhaps that message was already purged.
// TODO: provide a means for caller to test if sequence number is within range of stored history.
func (self *History) Recall(sequence SeqNum) []byte {
	message, ok := self.messages[sequence]
	if !ok { return nil }
	return message
}

func (self *History) purgeOldest() {

	// We'll always keep at least the most recent message sent.
	// The age and storage parameters actually don't consider the size or the true age of the most recent message.
	// The age of the oldest message is relative to the age of the newest, and the storage
	// used counts all but the size of the newset.

	for {
		prevLen := len(self.sequences)
		if prevLen <= 1 { break }

		oldestSeq := self.sequences[0]
		newestSeq := self.sequences[prevLen-1]

		age := AgeNanos(self.ticks[newestSeq] - self.ticks[oldestSeq])
		storage := SizeBytes(newestSeq - oldestSeq)

		if age < self.minAge { break }
		if storage < self.maxStorage && age < self.maxAge { break }

		self.sequences = self.sequences[1:prevLen]
		delete(self.messages, oldestSeq)
		delete(self.ticks, oldestSeq)
	}
}


