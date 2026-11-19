package bridge

import (
	"encoding/binary"
	"testing"
)

func testRTPPacket(seq uint16) []byte {
	p := make([]byte, 12)
	p[0] = 0x80
	p[1] = 127
	binary.BigEndian.PutUint16(p[2:4], seq)
	return p
}

func seqs(pkts [][]byte) []uint16 {
	out := make([]uint16, 0, len(pkts))
	for _, p := range pkts {
		out = append(out, binary.BigEndian.Uint16(p[2:4]))
	}
	return out
}

func TestRTPReorderBuffer(t *testing.T) {
	r := newRTPReorderBuffer(64)
	var got []uint16
	for _, s := range []uint16{0, 1, 3, 4, 5, 6, 7, 8, 2, 9} {
		got = append(got, seqs(r.Push(testRTPPacket(s)))...)
	}
	want := []uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
	if out := r.Push(testRTPPacket(3)); len(out) != 0 {
		t.Fatalf("late duplicate emitted: %v", seqs(out))
	}

	w := newRTPReorderBuffer(2)
	var skipped []uint16
	for _, s := range []uint16{10, 12, 13} {
		skipped = append(skipped, seqs(w.Push(testRTPPacket(s)))...)
	}
	wantSkipped := []uint16{10, 12, 13}
	for i := range wantSkipped {
		if i >= len(skipped) || skipped[i] != wantSkipped[i] {
			t.Fatalf("window skip got %v, want %v", skipped, wantSkipped)
		}
	}
}
