package cpu

import "testing"

func TestLda(t *testing.T) {
	// LDA $1080, X
	Memory[0x0000] = 0xBD
	Memory[0x0001] = 0x80
	Memory[0x0002] = 0x10

	x = 4

	Memory[0x1084] = 0x76

	pc = 0

	CodeToOp[Memory[pc]]()
	if a != Memory[0x1084] {
		t.Error("Expected ", Memory[0x1084], ", got ", a)
	}
}
