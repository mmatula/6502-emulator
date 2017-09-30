package cpu

import "go/ast"

// helper functions
func setZeroAndNegative(value uint8) {
	zero = value == 0
	negative = value >= 0x80
}

func readUInt16(address uint16) uint16 {
	return uint16(Memory[address]) + uint16(Memory[address + 1]) << 8
}

func pushByte(value uint8) {
	Memory[0x0100 + sp] = value
	sp--
}

func pushWord(value uint16) {
	pushByte(uint8(value >> 8))
	pushByte(uint8(value))
}

// base for branches
func jumpIfTrue(condition bool) {
	if condition {
		pc = uint16(int(pc) + int(int8(Memory[pc + 1])))
	} else {
		pc += 2
	}
}

// ------ VALUE ADDRESSING MODES -------
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
	return &Memory[absoluteAddress()]
}

func absoluteX() *uint8 {
	return &Memory[absoluteAddress() + uint16(x)]
}

func absoluteY() *uint8 {
	return &Memory[absoluteAddress() + uint16(y)]
}

func indexedIndirect() *uint8 {
	return &Memory[readUInt16(uint16(Memory[pc + 1] + x))]
}

func indirectIndexed() *uint8 {
	return &Memory[readUInt16(uint16(Memory[pc + 1])) + uint16(y)]
}

// ------ POINTER ADDRESSING MODES ------
func absoluteAddress() uint16 {
	return readUInt16(pc + 1)
}

func indirectAddress() uint16 {
	return readUInt16(absoluteAddress())
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

func brk() {
	breakCommand = true
	// TODO: dodělat
}

func bvc() {
	jumpIfTrue(!overflow)
}

func bvs() {
	jumpIfTrue(overflow)
}

func clc() {
	carry = false
	pc++
}

func cld() {
	decimalMode = false
	pc++
}

func cli() {
	interruptDisable = false
	pc++
}

func clv() {
	overflow = false
	pc++
}

func cmp(addressingMode func() *uint8) {
	value := *addressingMode()
	carry = a >= value
	setZeroAndNegative(value)
}

func cpx(addressingMode func() *uint8) {
	value := *addressingMode()
	carry = x >= value
	setZeroAndNegative(x)
}

func cpy(addressingMode func() *uint8) {
	value := *addressingMode()
	carry = y >= value
	setZeroAndNegative(y)
}

func dec(addressingMode func() *uint8) {
	ptr := addressingMode()
	*ptr--
	setZeroAndNegative(*ptr)
}

func dex() {
	x--
	setZeroAndNegative(x)
	pc++
}

func dey() {
	y--
	setZeroAndNegative(y)
	pc++
}

func eor(addressingMode func() *uint8) {
	a ^= *addressingMode()
	setZeroAndNegative(a)
}

func inc(addressingMode func() *uint8) {
	ptr := addressingMode()
	*ptr++
	setZeroAndNegative(*ptr)
}

func inx() {
	x++
	setZeroAndNegative(x)
	pc++
}

func iny() {
	y++
	setZeroAndNegative(y)
	pc++
}

func jmp(addressingMode func() uint16) {
	pc = addressingMode()
}

func jsr() {
	pushWord(pc + 2)
	pc = absoluteAddress()
}

func lsr(addressingMode func() *uint8) {
	ptr := addressingMode()
	carry = *ptr & 1 == 1
	*ptr >>= 1
	zero = *ptr == 0
	negative = false
}

func nop() {
	pc++
}