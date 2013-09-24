// header_test.go

package header

import (
	"testing"
)

func TestMakeMessageHeader(t *testing.T) {
	const seq = uint32(23)

	h := MakeMessageHeader(seq)

	if !h.Valid() {
		t.Fail()
	}

	if h.MessageType() != Message {
		t.Fail()
	}

	if h.GetSequence() != seq {
		t.Fail()
	}
}
