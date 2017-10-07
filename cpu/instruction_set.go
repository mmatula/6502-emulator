package cpu

import "errors"

type amode func(bool) *uint8

// helper functions
func setZeroAndNegative(value uint8) {
	zero = value == 0
	negative = value >= 0x80
}

func readUInt16(address uint16) uint16 {
	ticksToNext += 2
	return uint16(Memory[address]) + (uint16(Memory[address + 1]) << 8)
}

func readUInt16WithError(address uint16) uint16 {
	addressLo := address & 0x00FF
	addressHi := address & 0xFF00
	ticksToNext += 2
	return uint16(Memory[address]) + (uint16(Memory[addressHi + ((addressLo + 1) & 0x00FF)]) << 8)
}

func pushByte(value uint8) {
	Stack[sp] = value
	sp--
	ticksToNext++
}

func pushWord(value uint16) {
	pushByte(uint8(value >> 8))
	pushByte(uint8(value))
}

func pullByte() uint8 {
	sp++
	ticksToNext++
	return Stack[sp]
}

func pullWord() uint16 {
	addressLo := pullByte()
	addressHi := pullByte()
	return uint16(addressLo) + (uint16(addressHi) << 8)
}

func getBit(value, index uint8) bool {
	return value & (1 << index) != 0
}

// base for branches
func jumpIfTrue(condition, isBvc bool) {
	if condition {
		ticksToNext++
		oldPc := pc + 2
		pc = uint16(int(pc) + int(int8(Memory[pc + 1])))
		if !isBvc && (oldPc & 0xFF00 < pc & 0xFF00) {
			ticksToNext++
		}
	} else {
		pc += 2
	}
	ticksToNext++
}

func consumeTicksForWrite(write bool) {
	if write {
		ticksToNext += 2
	}
}

func incrementPc(increment uint16) {
	pc += increment
}

// ------ VALUE ADDRESSING MODES -------
func immediate(write bool) *uint8 {
	if write {
		errors.New("should never happen")
	}
	ticksToNext++
	defer incrementPc(2)
	return &Memory[pc + 1]
}

func accumulator(write bool) *uint8 {
	if !write {
		errors.New("should never happen")
	}
	ticksToNext++
	pc += 1
	return &a
}

func zeroPageIndexed(address uint16, index uint8, write bool) *uint8 {
	consumeTicksForWrite(write)
	pc += 2
	ticksToNext += 2
	return &Memory[Memory[address] + index]
}

func zeroPage(write bool) *uint8 {
	return zeroPageIndexed(pc + 1, 0, write)
}

func zeroPageX(write bool) *uint8 {
	ticksToNext++
	return zeroPageIndexed(pc + 1, x, write)
}

func zeroPageY(write bool) *uint8 {
	ticksToNext++
	return zeroPageIndexed(pc + 1, y, write)
}

func absolute(write bool) *uint8 {
	consumeTicksForWrite(write)
	defer incrementPc(3)
	ticksToNext++
	return &Memory[absoluteAddress()]
}

func absoluteIndexed(address uint16, index uint8, write bool) *uint8 {
	ticksToNext++
	indexedAddress := address + uint16(index)
	if write || (address & 0xFF00 < indexedAddress & 0xFF00) {
		ticksToNext++
	}
	pc += 3
	consumeTicksForWrite(write)
	return &Memory[indexedAddress]
}

func absoluteX(write bool) *uint8 {
	return absoluteIndexed(absoluteAddress(), x, write)
}

func absoluteY(write bool) *uint8 {
	return absoluteIndexed(absoluteAddress(), y, write)
}

func indexedIndirect(write bool) *uint8 {
	if write {
		errors.New("should never happen")
	}
	defer incrementPc(2)
	ticksToNext += 3
	return &Memory[readUInt16(uint16(Memory[pc + 1] + x))]
}

func indirectIndexed(write bool) *uint8 {
	if write {
		errors.New("should never happen")
	}
	defer incrementPc(2)
	ticksToNext++
	return absoluteIndexed(readUInt16(uint16(Memory[pc + 1])), y, write)
}

// ------ POINTER ADDRESSING MODES ------
func absoluteAddress() uint16 {
	return readUInt16(pc + 1)
}

func indirectAddress() uint16 {
	return readUInt16WithError(absoluteAddress())
}

// ------ INSTRUCTIONS -------
func adc(addressingMode amode) {
	value := *addressingMode(false)

	var carryInc uint8 = 0
	if carry {
		carryInc = 1
	}

	if decimalMode {
		aL := a & 0x0F + value & 0x0F + carryInc
		aH := a >> 4 + value >> 4
		if aL > 0x0F {
			aH++
		}
		if aL > 0x09 {
			aL += 0x06
		}
		zero = a + value + carryInc == 0
		negative = aH & 0x08 != 0
		overflow = (a & 0x80 == value & 0x80) && (value & 0x80 != (aH << 4) & 0x80)
		if aH > 0x09 {
			aH += 0x06
		}
		carry = aH > 0x0F
		a = (aH << 4) + (aL & 0x0F)
	} else {
		newValue := uint16(a) + uint16(value)
		if carry {
			newValue++
		}
		carry = newValue > 0xFF
		overflow = (a & 0x80 == value & 0x80) && (value & 0x80 != uint8(newValue) & 0x80)
		a = uint8(newValue)
		setZeroAndNegative(a)
	}
}

func and(addressingMode amode) {
	a &= *addressingMode(false)
	setZeroAndNegative(a)
}

func asl(addressingMode amode) {
	ptr := addressingMode(true)
	carry = *ptr >= 128
	*ptr <<= 1
	setZeroAndNegative(*ptr)
}

func bcc() {
	jumpIfTrue(!carry, false)
}

func bcs() {
	jumpIfTrue(carry, false)
}

func beq() {
	jumpIfTrue(zero, false)
}

func bit(addressingMode amode) {
	value := *addressingMode(false)
	negative = value & 128 == 128
	overflow = value & 64 == 64
	zero = value & a == 0
}

func bmi() {
	jumpIfTrue(negative, false)
}

func bne() {
	jumpIfTrue(!zero, false)
}

func bpl() {
	jumpIfTrue(!negative, false)
}

func brk() {
	pushWord(pc + 2)
	pushByte(getPs())
	pc = 0xFFFE
	ticksToNext += 3
}

func bvc() {
	jumpIfTrue(!overflow, true)
}

func bvs() {
	jumpIfTrue(overflow, false)
}

func clc() {
	carry = false
	pc++
	ticksToNext++
}

func cld() {
	decimalMode = false
	pc++
	ticksToNext++
}

func cli() {
	interruptDisable = false
	pc++
	ticksToNext++
}

func clv() {
	overflow = false
	pc++
	ticksToNext++
}

func cmp(addressingMode amode) {
	value := *addressingMode(false)
	carry = a >= value
	setZeroAndNegative(value)
}

func cpx(addressingMode amode) {
	value := *addressingMode(false)
	carry = x >= value
	setZeroAndNegative(x)
}

func cpy(addressingMode amode) {
	value := *addressingMode(false)
	carry = y >= value
	setZeroAndNegative(y)
}

func dec(addressingMode amode) {
	ptr := addressingMode(true)
	*ptr--
	setZeroAndNegative(*ptr)
}

func dex() {
	x--
	setZeroAndNegative(x)
	pc++
	ticksToNext++
}

func dey() {
	y--
	setZeroAndNegative(y)
	pc++
	ticksToNext++
}

func eor(addressingMode amode) {
	a ^= *addressingMode(false)
	setZeroAndNegative(a)
}

func inc(addressingMode amode) {
	ptr := addressingMode(true)
	*ptr++
	setZeroAndNegative(*ptr)
}

func inx() {
	x++
	setZeroAndNegative(x)
	pc++
	ticksToNext++
}

func iny() {
	y++
	setZeroAndNegative(y)
	pc++
	ticksToNext++
}

func jmp(addressingMode func() uint16) {
	pc = addressingMode()
}

func jsr() {
	pushWord(pc + 2)
	pc = absoluteAddress()
	ticksToNext++
}

func lda(addressingMode amode) {
	a = *addressingMode(false)
	setZeroAndNegative(a)
}

func ldx(addressingMode amode) {
	x = *addressingMode(false)
	setZeroAndNegative(x)
}

func ldy(addressingMode amode) {
	y = *addressingMode(false)
	setZeroAndNegative(y)
}

func lsr(addressingMode amode) {
	ptr := addressingMode(true)
	carry = *ptr & 1 == 1
	*ptr >>= 1
	zero = *ptr == 0
	negative = false
}

func nop() {
	pc++
	ticksToNext++
}

func ora(addressingMode amode) {
	a |= *addressingMode(false)
	setZeroAndNegative(a)
}

func pha() {
	pushByte(a)
	pc++
	ticksToNext++
}

func php() {
	pushByte(getPs())
	pc++
	ticksToNext++
}

func pla() {
	a = pullByte()
	pc++
	ticksToNext += 2
}

func plp() {
	setPs(pullByte())
	pc++
	ticksToNext += 2
}

func rol(addressingMode amode) {
	ptr := addressingMode(true)
	newCarry := getBit(*ptr, 7)
	*ptr <<= 1
	if carry {
		*ptr++
	}
	carry = newCarry
	setZeroAndNegative(*ptr)
}

func ror(addressingMode amode) {
	ptr := addressingMode(true)
	newCarry := getBit(*ptr, 0)
	*ptr >>= 1
	if carry {
		*ptr += 128
	}
	carry = newCarry
	setZeroAndNegative(*ptr)
}

func rti() {
	setPs(pullByte())
	pc = pullWord()
	ticksToNext += 2
}

func rts() {
	pc = pullWord() + 1
	ticksToNext += 3
}

func sbc(addressingMode amode) {
	value := *addressingMode(false)

	var carryDec uint8 = 0
	if !carry {
		carryDec = 1
	}

	newValue := uint16(a) - uint16(value) - uint16(carryDec)
	overflow = (a & 0x80 != value & 0x80) && (value & 0x80 == uint8(newValue) & 0x80)

	if decimalMode {
		aL := a & 0x0F - value & 0x0F - carryDec
		aH := a >> 4 - value >> 4
		if aL > 0x0F {
			aL -= 0x06
			aH--
		}
		if aH > 0x0F {
			aH -= 0x06
		}
		a = (aH << 4) + (aL & 0x0F)
	} else {
		a = uint8(newValue)
	}
	carry = newValue < 0x0100
	setZeroAndNegative(a)
}

func sec() {
	carry = true
	pc++
	ticksToNext++
}

func sed() {
	decimalMode = true
	pc++
	ticksToNext++
}

func sei() {
	interruptDisable = true
	pc++
	ticksToNext++
}

func sta(addressingMode amode) {
	*addressingMode(true) = a
	ticksToNext -= 2
}

func stx(addressingMode amode) {
	*addressingMode(true) = x
	ticksToNext -= 2
}

func sty(addressingMode amode) {
	*addressingMode(true) = y
	ticksToNext -= 2
}

func tax() {
	x = a
	setZeroAndNegative(x)
	pc++
	ticksToNext++
}

func tay() {
	y = a
	setZeroAndNegative(y)
	pc++
	ticksToNext++
}

func tsx() {
	x = sp
	setZeroAndNegative(x)
	pc++
	ticksToNext++
}

func txa() {
	a = x
	setZeroAndNegative(a)
	pc++
	ticksToNext++
}

func txs() {
	sp = x
	pc++
	ticksToNext++
}

func tya() {
	a = y
	setZeroAndNegative(a)
	pc++
	ticksToNext++
}
