package main

// registers
var pc uint16	// program counter
var sp uint8	// stack pointer
var a uint8		// accumulator
var x, y uint8	// index registers
var ps uint8	// processor status (flags)

// memory
var memory [65536]uint8

func main() {
	// LDA $23
	memory[0x0000] = 0xA9
	memory[0x0001] = 0x23
	memory[0x0002] = 0x00
	memory[0x0003] = 0x00
	memory[0x0004] = 0x00
	memory[0x0005] = 0x00

	//...
}

func opLdaImmediate() {
	// LDA $23
	// op code = $A9
	a = memory[pc + 1]
	pc = pc + 2
}

func opLdaZeroPage() {
	// LDA [$80]
	// op code = $A5
	var address = memory[pc + 1]
	a = memory[address]
	pc = pc + 2
}