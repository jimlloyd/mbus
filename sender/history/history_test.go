// mbus/sender/history/history_test.go

package history

import (
	"testing"
)

func TestNominal(t *testing.T) {

	minAgeSeconds := int32(0)
	maxAgeSeconds := int32(10)
	maxPayloadMB := uint32(1)

	hist := NewHistory(minAgeSeconds, maxAgeSeconds, maxPayloadMB)
	if hist.Length() != 0 {
		t.Error("New History should have length 0")
	}
	if hist.Recall(0) != nil {
		t.Error("Recall on empty history should return nil")
	}

	s1 := "message1"
	s2 := "message2"
	s3 := "message3"

	hist.Add(0, []byte(s1))
	if hist.Length() != 1 {
		t.Error("New History should have length 1")
	}
	if string(hist.Recall(0)) != s1 {
		t.Error("Recall failed to return correct message")
	}

	hist.Add(10, []byte(s2))
	if hist.Length() != 2 {
		t.Error("New History should have length 2")
	}
	if string(hist.Recall(10)) != s2 {
		t.Error("Recall failed to return correct message")
	}

	// The large sequence number here should trigger a purge.
	hist.Add(1e8, []byte(s3))
	if hist.Length() != 1 {
		t.Error("New History should have length 1")
	}
	if hist.Recall(0) != nil {
		t.Error("Recall failed to return nil for purged message")
	}
	if string(hist.Recall(1e8)) != s3 {
		t.Error("Recall failed to return correct message")
	}
}

