package asnrt

import "testing"

// P4-14: named-bit value accessors mirroring asnjc value.NamedBits and
// value.BitString.getBit(...)/setBit(...) (NamedNumber variants).

func TestBitStringNamedBitHelpers(t *testing.T) {
	bs := BitString{}
	// Name->number mapping for a type whose named bits are: alpha=0, beta=4, gamma=7.
	alpha := NewNamedNumber("alpha", 0)
	beta := NewNamedNumber("beta", 4)
	gamma := NewNamedNumber("gamma", 7)

	if bs.GetNamedBit(alpha) {
		t.Fatalf("alpha should start cleared")
	}

	bs.SetNamedBit(alpha, true)
	bs.SetNamedBit(gamma, true)

	if !bs.GetNamedBit(alpha) {
		t.Fatalf("alpha should be set")
	}
	if bs.GetNamedBit(beta) {
		t.Fatalf("beta should be cleared")
	}
	if !bs.GetNamedBit(gamma) {
		t.Fatalf("gamma should be set")
	}

	bs.SetNamedBit(alpha, false)
	if bs.GetNamedBit(alpha) {
		t.Fatalf("alpha should be cleared after SetNamedBit(false)")
	}

	if bs.Length() != 8 {
		t.Fatalf("expected Length()==8 (bit 7 set), got %d", bs.Length())
	}

	// Cross-check against the raw byte accessor.
	if bs.Bytes[0] == 0 {
		t.Fatalf("expected byte 0 to be non-zero (bit 7 set)")
	}
}

func TestNamedBitsMap(t *testing.T) {
	flag := NewNamedNumber("flag", 3)
	nb := NewNamedBits(nil)

	if nb.GetNamedBit(flag) {
		t.Fatalf("unrecorded named bit should read false")
	}
	// nil map must not panic on read.
	if nb.GetBit(3) {
		t.Fatalf("unrecorded bit 3 should read false")
	}

	nb = NewNamedBits(map[int]bool{})
	nb.SetNamedBit(flag, true)
	if !nb.GetNamedBit(flag) {
		t.Fatalf("flag should be set")
	}
	nb.SetNamedBit(flag, false)
	if nb.GetNamedBit(flag) {
		t.Fatalf("flag should be cleared after SetNamedBit(false)")
	}

	// Distinguish "recorded+cleared" from "absent": absent reads false, recorded
	// cleared also reads false, recorded set reads true.
	absent := NewNamedNumber("absent", 9)
	if nb.GetNamedBit(absent) {
		t.Fatalf("absent named bit should read false")
	}
}