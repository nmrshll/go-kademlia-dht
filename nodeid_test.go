package dht

import (
	"testing"
)

func TestNodeID(t *testing.T) {
	a := NodeID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	aBis := NodeID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	b := NodeID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 19, 18}
	c := NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1}

	if !a.Equals(a) {
		t.Fatalf("nodeID 'a' should be equal to itself")
	}
	if !a.Equals(aBis) {
		t.Fatalf("nodeID 'a' should be equal to aBis")
	}
	if a.Equals(b) {
		t.Fatalf("nodeID 'a' should not be equal to nodeID 'b'")
	}

	if !a.Xor(b).Equals(c) {
		t.Error(a.Xor(b))
	}

	if a.LeadingZeros() != 15 {
		t.Fatalf("expected a to have length 15")
	}
	if b.LeadingZeros() != 15 {
		t.Fatalf("expected b to have length 15")
	}
	if c.LeadingZeros() != 151 {
		t.Fatalf("expected c to be length 151")
	}

	if b.Less(a) {
		t.Fatalf("expected b to be bigger than a")
	}
	if !a.Less(b) {
		t.Fatalf("expected b to be bigger than a")
	}

	strID := "0123456789abcdef0123456789abcdef01234567"
	if NewNodeID(strID).String() != strID {
		t.Error(NewNodeID(strID).String())
	}
}
