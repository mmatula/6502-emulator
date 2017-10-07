package cpu

// helper functions
func setZeroAndNegative(value uint8) {
	zero = value == 0
	negative = value >= 0x80
}

func readUInt16(address uint16) uint16 {
	return uint16(Memory[address]) + (uint16(Memory[address + 1]) << 8)
}

func readUInt16WithError(address uint16) uint16 {
	addressLo := address & 0x00FF
	addressHi := address & 0xFF00
	return uint16(Memory[address]) + (uint16(Memory[addressHi + ((addressLo + 1) & 0x00FF)]) << 8)
}

func pushByte(value uint8) {
	Stack[sp] = value
	sp--
}

func pushWord(value uint16) {
	pushByte(uint8(value >> 8))
	pushByte(uint8(value))
}

func pullByte() uint8 {
	sp++
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
	return readUInt16WithError(absoluteAddress())
}

// ------ INSTRUCTIONS -------
func adc(addressingMode func() *uint8) {
	value := *addressingMode()

	var carryInc uint8 = 0
	if carry {
		carryInc = 1
	}

	if decimalMode {
		aL := a & 0x0F + value & 0x0F + carryInc;
		aH := a >> 4 + value >> 4;
		if aL > 0x0F {
			aH++
		}
		if aL > 0x09 {
			aL += 0x06
		}
		zero = a + value + carryInc == 0
		negative = aH & 0x08 != 0
		overflow = (a & 0x80 == value & 0x80) && (a & 0x80 != (aH << 4) & 0x80)
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
		overflow = (a & 0x80 == value & 0x80) && (a & 0x80 != uint8(newValue) & 0x80)
		a = uint8(newValue)
		setZeroAndNegative(a)
	}
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
	pushWord(pc + 2)
	pushByte(getPs())
	pc = 0xFFFE
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

func ora(addressingMode func() *uint8) {
	a |= *addressingMode()
	setZeroAndNegative(a)
}

func pha() {
	pushByte(a)
	pc++
}

func php() {
	pushByte(getPs())
	pc++
}

func pla() {
	a = pullByte()
	pc++
}

func plp() {
	setPs(pullByte())
	pc++
}

func rol(addressingMode func() *uint8) {
	ptr := addressingMode()
	newCarry := getBit(*ptr, 7)
	*ptr <<= 1
	if carry {
		*ptr++
	}
	carry = newCarry
	setZeroAndNegative(*ptr)
}

func ror(addressingMode func() *uint8) {
	ptr := addressingMode()
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
}

func rts() {
	pc = pullWord() + 1
}

func sbc(addressingMode func() *uint8) {
	// TODO: implement
}

func sec() {
	carry = true
	pc++
}

func sed() {
	decimalMode = true
	pc++
}

func sei() {
	interruptDisable = true
	pc++
}

func sta(addressingMode func() *uint8) {
	*addressingMode() = a
}

func stx(addressingMode func() *uint8) {
	*addressingMode() = x
}

func sty(addressingMode func() *uint8) {
	*addressingMode() = y
}

func tax() {
	x = a
	setZeroAndNegative(x)
	pc++
}

func tay() {
	y = a
	setZeroAndNegative(y)
	pc++
}

func tsx() {
	x = sp
	setZeroAndNegative(x)
	pc++
}

func txa() {
	a = x
	setZeroAndNegative(a)
	pc++
}

func txs() {
	sp = x
	pc++
}

func tya() {
	a = y
	setZeroAndNegative(a)
	pc++
}
