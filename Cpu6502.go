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
	memory[0x0000] = 0xA5
	memory[0x0001] = 0x00
	memory[0x0002] = 0x00
	memory[0x0003] = 0x00
	memory[0x0004] = 0x00
	memory[0x0005] = 0x00

	//...
}

func opLdaImmediate() {
	// LDA #23
	// opcode = $A9
	a = memory[pc + 1]
	pc += 2
}

func opLdaZeroPage() {
	// LDA $80
	// opcode = $A5
	var address = memory[pc + 1]
	a = memory[address]
	pc += 2
}

func opLdaZeroPageX() {
	// LDA $80, X
	// opcode = $B5
	var address = memory[pc + 1]
	a = memory[address + x]
	pc += 2
}

func opLdaAbsolute() {
	// LDA $1080
	// opcode = $AD
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a = memory[address]
	pc += 3
}

func opLdaAbsoluteX() {
	// LDA $1080, X
	// opcode = $BD
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a = memory[address + uint16(x)]
	pc += 3
}

func opLdaAbsoluteY() {
	// LDA $1080, Y
	// opcode = $B9
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a = memory[address + uint16(y)]
	pc += 3
}

func opLdaIndexedIndirect() {
	// LDA ($80, X)
	// opcode = $A1
	var addressOfAddress = memory[pc + 1]
	var addressLo = memory[addressOfAddress + x]
	var addressHi = memory[addressOfAddress + x + 1]
	var address = uint16(addressLo) + uint16(addressHi) * 256
	a = memory[address]
	pc += 2
}

func opLdaIndirectIdexed() {
	// LDA ($80), Y
	// opcode = $B1
	var addressOfAddress = memory[pc + 1]
	var addressLo = memory[addressOfAddress]
	var addressHi = memory[addressOfAddress + 1]
	var address = uint16(addressLo) + uint16(addressHi) * 256
	a = memory[address + uint16(y)]
	pc += 2
}