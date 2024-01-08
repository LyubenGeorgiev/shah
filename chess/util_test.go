package chess

import "testing"

func TestSetGetBit(t *testing.T) {

	for i := square(0); i < 64; i++ {
		b := Bitboard(0)
		b.setBit(i)
		for j := square(0); j < 64; j++ {
			if j != i && b.getBit(j) {
				t.Fatalf("Bit %d should be zero but it is one\n", int(j))
			} else if j == i && !b.getBit(j) {
				t.Fatalf("Bit %d should be one but it is zero\n", int(j))
			}
		}
	}
}
