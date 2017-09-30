package cpu

// helper functions
func setZeroAndNegative(value uint8) {
	zero = value == 0
	negative = value >= 0x80
}

func readUInt16(address uint16) uint16 {
	return uint16(Memory[address]) + uint16(Memory[address + 1]) << 8
}

// base for branches
func jumpIfTrue(condition bool) {
	if condition {
		pc = uint16(int(pc) + int(int8(Memory[pc + 1])))
	} else {
		pc += 2
	}
}

// ------ ADDRESSING MODES -------
func immediate() *uint8 {
	return &Memory[pc + 1]
}

func accumulator() *uint8 {
	return &a
}

func zeroPage() *uint8 {
	return &Memory[Memory[pc + 1]]
}

func zeroPageX() *uint8 {
	return &Memory[Memory[pc + 1] + x]
}

func zeroPageY() *uint8 {
	return &Memory[Memory[pc + 1] + y]
}

func absolute() *uint8 {
	return &Memory[readUInt16(pc + 1)]
}

func absoluteX() *uint8 {
	return &Memory[readUInt16(pc + 1) + uint16(x)]
}

func absoluteY() *uint8 {
	return &Memory[readUInt16(pc + 1) + uint16(y)]
}

func indexedIndirect() *uint8 {
	return &Memory[readUInt16(uint16(Memory[pc + 1] + x))]
}

func indirectIndexed() *uint8 {
	return &Memory[readUInt16(uint16(Memory[pc + 1])) + uint16(y)]
}

// ------ INSTRUCTIONS -------
func adc(addressingMode func() *uint8) {
	prevA := a
	value := *addressingMode()
	a += value
	if carry {
		a++
	}
	setZeroAndNegative(a)
	carry = a <= value
	overflow = (prevA >= 128 && a < 128) || (prevA < 128 && a >= 128)
}

func and(addressingMode func() *uint8) {
	a &= *addressingMode()
	setZeroAndNegative(a)
}

func asl(addressingMode func() *uint8) {
	ptr := addressingMode()
	carry = *ptr >= 128
	*ptr <<= 1
	setZeroAndNegative(*ptr)
}

func bcc() {
	jumpIfTrue(!carry)
}

func bcs() {
	jumpIfTrue(carry)
}

func beq() {
	jumpIfTrue(zero)
}

func bit(addressingMode func() *uint8) {
	value := *addressingMode()
	negative = value & 128 == 128
	overflow = value & 64 == 64
	zero = value & a == 0
}

func bmi() {
	jumpIfTrue(negative)
}

func bne() {
	jumpIfTrue(!zero)
}

func bpl() {
	jumpIfTrue(!negative)
}

func lda(addressingMode func() *uint8) {
	a = *addressingMode()
	setZeroAndNegative(a)
}

func ldx(addressingMode func() *uint8) {
	x = *addressingMode()
	setZeroAndNegative(x)
}

func ldy(addressingMode func() *uint8) {
	y = *addressingMode()
	setZeroAndNegative(y)
}
