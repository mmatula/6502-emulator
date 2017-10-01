package cpu

import "testing"

func TestAdcDecimal(t *testing.T) {
	decimalMode = true
	a = 0x79
	carry = false
	pc = 0

	// ADC #28
	Memory[0x0000] = 0x69
	Memory[0x0001] = 0x28

	CodeToOp[Memory[pc]]()

	if a != 0x07 {
		t.Error("Expected 7, got ", a)
	}
	if !carry {
		t.Error("Expected carry flag to be set")
	}
}

func TestAdcZeroDecimal(t *testing.T) {
	decimalMode = true
	a = 0
	carry = false
	pc = 0

	// ADC #0
	Memory[0x0000] = 0x69
	Memory[0x0001] = 0x00

	CodeToOp[Memory[pc]]()

	if a != 0 {
		t.Error("Expected 0, got ", a)
	}
	if !zero {
		t.Error("Expected zero flag to be set")
	}
}

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

func TestJmpIndirectPageCross(t *testing.T) {
	// JMP ($02FF)
	Memory[0x0000] = 0x6C
	Memory[0x0001] = 0xFF
	Memory[0x0002] = 0x02

	Memory[0x0200] = 0x3C
	Memory[0x02FF] = 0xFF
	Memory[0x0300] = 0xFF

	pc = 0
	CodeToOp[Memory[pc]]()
	if pc != 0x3CFF {
		t.Error("Expected ", 0x3CFF, ", got ", pc)
	}
}