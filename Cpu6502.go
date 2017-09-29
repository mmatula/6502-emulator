package main

// registers
var pc uint16	// program counter
var sp uint8	// stack pointer
var a uint8		// accumulator
var x, y uint8	// index registers
var carry, zero, interruptDisable, decimalMode, breakCommand, overflow, negative bool	// processor status (flags)

// memory
var memory [65536]uint8

func setZeroAndNegative(value uint8) {
	zero = value == 0
	negative = value >= 0x80
}
func addToAccumulator(value uint8) {
	prevA := a
	a += value
	if carry {
		a++
	}
	setZeroAndNegative(a)
	carry = a <= value
	overflow = (prevA >= 128 && a < 128) || (prevA < 128 && a >= 128)
}
func jumpIfTrue(condition bool) {
	if condition {
		pc = uint16(int(pc) + int(int8(memory[pc + 1])))
	} else {
		pc += 2
	}
}
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


// ADC:
func opAdcImmediate() {
	// ADC #23
	// opcode = $69
	addToAccumulator(memory[pc + 1])
	pc += 2
}

func opAdcZeroPage() {
	// ADC $80
	// opcode = $65
	var address = memory[pc + 1]
	addToAccumulator(memory[address])
	pc += 2
}

func opAdcZeroPageX() {
	// ADC $80, X
	// opcode = $75
	var address = memory[pc + 1]
	addToAccumulator(memory[address + x])
	pc += 2
}

func opAdcAbsolute() {
	// ADC $1080
	// opcode = $6D
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	addToAccumulator(memory[address])
	pc += 3
}

func opAdcAbsoluteX() {
	// ADC $1080, X
	// opcode = $7D
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	addToAccumulator(memory[address + uint16(x)])
	pc += 3
}

func opAdcAbsoluteY() {
	// ADC $1080, Y
	// opcode = $79
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	addToAccumulator(memory[address + uint16(y)])
	pc += 3
}

func opAdcIndexedIndirect() {
	// ADC ($80, X)
	// opcode = $61
	var addressOfAddress = memory[pc + 1]
	var addressLo = memory[addressOfAddress + x]
	var addressHi = memory[addressOfAddress + x + 1]
	var address = uint16(addressLo) + uint16(addressHi) * 256
	addToAccumulator(memory[address])
	pc += 2
}

func opAdcIndirectIndexed() {
	// ADC ($80), Y
	// opcode = $71
	var addressOfAddress = memory[pc + 1]
	var addressLo = memory[addressOfAddress]
	var addressHi = memory[addressOfAddress + 1]
	var address = uint16(addressLo) + uint16(addressHi) * 256
	addToAccumulator(memory[address + uint16(y)])
	pc += 2
}


// AND:
func opAndImmediate() {
	// AND #23
	// opcode = $29
	a &= memory[pc + 1]
	setZeroAndNegative(a)
	pc += 2
}

func opAndZeroPage() {
	// AND $80
	// opcode = $25
	var address = memory[pc + 1]
	a &= memory[address]
	setZeroAndNegative(a)
	pc += 2
}

func opAndZeroPageX() {
	// AND $80, X
	// opcode = $35
	var address = memory[pc + 1]
	a &= memory[address + x]
	setZeroAndNegative(a)
	pc += 2
}

func opAndAbsolute() {
	// AND $1080
	// opcode = $2D
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a &= memory[address]
	setZeroAndNegative(a)
	pc += 3
}

func opAndAbsoluteX() {
	// AND $1080, X
	// opcode = $3D
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a &= memory[address + uint16(x)]
	setZeroAndNegative(a)
	pc += 3
}

func opAndAbsoluteY() {
	// AND $1080, Y
	// opcode = $39
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a &= memory[address + uint16(y)]
	setZeroAndNegative(a)
	pc += 3
}

func opAndIndexedIndirect() {
	// AND ($80, X)
	// opcode = $21
	var addressOfAddress = memory[pc + 1]
	var addressLo = memory[addressOfAddress + x]
	var addressHi = memory[addressOfAddress + x + 1]
	var address = uint16(addressLo) + uint16(addressHi) * 256
	a &= memory[address]
	setZeroAndNegative(a)
	pc += 2
}

func opAndIndirectIndexed() {
	// AND ($80), Y
	// opcode = $31
	var addressOfAddress = memory[pc + 1]
	var addressLo = memory[addressOfAddress]
	var addressHi = memory[addressOfAddress + 1]
	var address = uint16(addressLo) + uint16(addressHi) * 256
	a &= memory[address + uint16(y)]
	setZeroAndNegative(a)
	pc += 2
}


// ASL:
func shiftLeft(value uint8) uint8 {
	carry = value >= 128
	result := value << 1
	setZeroAndNegative(result)
	return result
}

func opAslAccumulator() {
	// ASL A
	// opcode = $0A
	a = shiftLeft(a)
	pc += 1
}

func opAslZeroPage() {
	// ASL $80
	// opcode = $06
	var address = memory[pc + 1]
	memory[address] = shiftLeft(memory[address])
	pc += 2
}

func opAslZeroPageX() {
	// ASL $80, X
	// opcode = $16
	var address = memory[pc + 1]
	memory[address + x] = shiftLeft(memory[address + x])
	pc += 2
}

func opAslAbsolute() {
	// ASL $1080
	// opcode = $0E
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	memory[address] = shiftLeft(memory[address])
	pc += 3
}

func opAslAbsoluteX() {
	// ASL $1080, X
	// opcode = $1E
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	memory[address + uint16(x)] = shiftLeft(memory[address + uint16(x)])
	setZeroAndNegative(a)
	pc += 3
}


// BCC:
func opBcc() {
	// BCC *+23
	// opcode = $90
	jumpIfTrue(!carry)
}


// BCS:
func opBcs() {
	// BCS *+23
	// opcode = $B0
	jumpIfTrue(carry)
}


// BEQ:
func opBeq() {
	// BEQ *+23
	// opcode = $F0
	jumpIfTrue(zero)
}


// BIT:
func opBitZeroPage() {
	// BIT $80
	// opcode = $24
	var address = memory[pc + 1]
	var value = memory[address]
	negative = value & 128 == 128
	overflow = value & 64 == 64
	zero = value & a == 0
	pc += 2
}

func opBitAbsolute() {
	// BIT $1080
	// opcode = $2C
	var address = uint16(memory[pc + 1]) + (uint16(memory[pc + 2]) << 8)
	var value = memory[address]
	negative = value & 128 == 128
	overflow = value & 64 == 64
	zero = value & a == 0
	pc += 3
}


// BMI:
func opBmi() {
	// BMI *+23
	// opcode = $30
	jumpIfTrue(negative)
}


// BNE:
func opBne() {
	// BNE *+23
	// opcode = $D0
	jumpIfTrue(!zero)
}


// BPL:
func opBpl() {
	// BPL *+23
	// opcode = $10
	jumpIfTrue(!negative)
}


// LDA:
func opLdaImmediate() {
	// LDA #23
	// opcode = $A9
	a = memory[pc + 1]
	setZeroAndNegative(a)
	pc += 2
}

func opLdaZeroPage() {
	// LDA $80
	// opcode = $A5
	var address = memory[pc + 1]
	a = memory[address]
	setZeroAndNegative(a)
	pc += 2
}

func opLdaZeroPageX() {
	// LDA $80, X
	// opcode = $B5
	var address = memory[pc + 1]
	a = memory[address + x]
	setZeroAndNegative(a)
	pc += 2
}

func opLdaAbsolute() {
	// LDA $1080
	// opcode = $AD
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a = memory[address]
	setZeroAndNegative(a)
	pc += 3
}

func opLdaAbsoluteX() {
	// LDA $1080, X
	// opcode = $BD
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a = memory[address + uint16(x)]
	setZeroAndNegative(a)
	pc += 3
}

func opLdaAbsoluteY() {
	// LDA $1080, Y
	// opcode = $B9
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	a = memory[address + uint16(y)]
	setZeroAndNegative(a)
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
	setZeroAndNegative(a)
	pc += 2
}

func opLdaIndirectIndexed() {
	// LDA ($80), Y
	// opcode = $B1
	var addressOfAddress = memory[pc + 1]
	var addressLo = memory[addressOfAddress]
	var addressHi = memory[addressOfAddress + 1]
	var address = uint16(addressLo) + uint16(addressHi) * 256
	a = memory[address + uint16(y)]
	setZeroAndNegative(a)
	pc += 2
}


// LDX:
func opLdxImmediate() {
	// LDX #23
	// opcode = $A2
	x = memory[pc + 1]
	setZeroAndNegative(x)
	pc += 2
}

func opLdxZeroPage() {
	// LDX $80
	// opcode = $A6
	var address = memory[pc + 1]
	x = memory[address]
	setZeroAndNegative(x)
	pc += 2
}

func opLdxZeroPageY() {
	// LDX $80, Y
	// opcode = $B6
	var address = memory[pc + 1]
	x = memory[address + y]
	setZeroAndNegative(x)
	pc += 2
}

func opLdxAbsolute() {
	// LDX $1080
	// opcode = $AE
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	x = memory[address]
	setZeroAndNegative(x)
	pc += 3
}

func opLdxAbsoluteY() {
	// LDX $1080, Y
	// opcode = $BE
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	x = memory[address + uint16(y)]
	setZeroAndNegative(x)
	pc += 3
}


// LDY:
func opLdyImmediate() {
	// LDY #23
	// opcode = $A0
	y = memory[pc + 1]
	setZeroAndNegative(y)
	pc += 2
}

func opLdyZeroPage() {
	// LDY $80
	// opcode = $A4
	var address = memory[pc + 1]
	y = memory[address]
	setZeroAndNegative(y)
	pc += 2
}

func opLdyZeroPageX() {
	// LDY $80, X
	// opcode = $B4
	var address = memory[pc + 1]
	y = memory[address + x]
	setZeroAndNegative(y)
	pc += 2
}

func opLdyAbsolute() {
	// LDY $1080
	// opcode = $AC
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	y = memory[address]
	setZeroAndNegative(y)
	pc += 3
}

func opLdyAbsoluteX() {
	// LDY $1080, X
	// opcode = $BC
	var address = uint16(memory[pc + 1]) + uint16(memory[pc + 2]) * 256
	y = memory[address + uint16(x)]
	setZeroAndNegative(y)
	pc += 3
}
