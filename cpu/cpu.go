package cpu

// registers
var pc uint16 = 0		// program counter
var sp uint8 = 0xFF		// stack pointer
var a uint8				// accumulator
var x, y uint8			// index registers
var carry, zero, interruptDisable, decimalMode, breakCommand, overflow, negative bool	// processor status (flags)

func getPs() uint8 {
	var result uint8 = 32
	if carry {
		result += 1
	}
	if zero {
		result += 2
	}
	if interruptDisable {
		result += 4
	}
	if decimalMode {
		result += 8
	}
	if breakCommand {
		result += 16
	}
	if overflow {
		result += 64
	}
	if negative {
		result += 128
	}
	return result
}

func setPs(ps uint8) {
	carry = ps & 1 == 1
	zero = ps & 2 == 2
	interruptDisable = ps & 4 == 4
	decimalMode = ps & 8 == 8
	breakCommand = ps & 16 == 16
	overflow = ps & 64 == 64
	negative = ps & 128 == 128
}

// memory
var Memory [65536]uint8
var Stack []uint8 = Memory[0x0100:0x01FF]

// instruction set
type op func()

var CodeToOp = [256]op {
	// 0x00 - 0x0F
	brk,
	nil,
	nil,
	nil,
	nil,
	nil,
	func() { asl(zeroPage); pc += 2 },
	nil,
	nil,
	nil,
	func() { asl(accumulator); pc++ },
	nil,
	nil,
	nil,
	func() { asl(absolute); pc += 3 },
	nil,

	// 0x10 - 0x1F
	bpl,
	nil,
	nil,
	nil,
	nil,
	nil,
	func() { asl(zeroPageX); pc += 2 },
	nil,
	clc,
	nil,
	nil,
	nil,
	nil,
	nil,
	func() { asl(absoluteX); pc += 3 },
	nil,

	// 0x20 - 0x2F
	jsr,
	func() { and(indexedIndirect); pc += 2 },
	nil,
	nil,
	func() { bit(zeroPage); pc += 2 },
	func() { and(zeroPage); pc += 2 },
	nil,
	nil,
	nil,
	func() { and(immediate); pc += 2 },
	nil,
	nil,
	func() { bit(absolute); pc += 3 },
	func() { and(absolute); pc += 3 },
	nil,
	nil,

	// 0x30 - 0x3F
	bmi,
	func() { and(indirectIndexed); pc += 2 },
	nil,
	nil,
	nil,
	func() { and(zeroPageX); pc += 2 },
	nil,
	nil,
	nil,
	func() { and(absoluteY); pc += 3 },
	nil,
	nil,
	nil,
	func() { and(absoluteX); pc += 3 },
	nil,
	nil,

	// 0x40 - 0x4F
	nil,
	func() { eor(indexedIndirect); pc += 2 },
	nil,
	nil,
	nil,
	func() { eor(zeroPage); pc += 2 },
	func() { lsr(zeroPage); pc += 2 },
	nil,
	nil,
	func() { eor(immediate); pc += 2 },
	func() { lsr(accumulator); pc++ },
	nil,
	func() { jmp(absoluteAddress) },
	func() { eor(absolute); pc += 3 },
	func() { lsr(absolute); pc += 3 },
	nil,

	// 0x50 - 0x5F
	bvc,
	func() { eor(indirectIndexed); pc += 2 },
	nil,
	nil,
	nil,
	func() { eor(zeroPageX); pc += 2 },
	func() { lsr(zeroPageX); pc += 2 },
	nil,
	cli,
	func() { eor(absoluteY); pc += 3 },
	nil,
	nil,
	nil,
	func() { eor(absoluteX); pc += 3 },
	func() { lsr(absoluteX); pc += 3 },
	nil,

	// 0x60 - 0x6F
	nil,
	func() { adc(indexedIndirect); pc += 2 },
	nil,
	nil,
	nil,
	func() { adc(zeroPage); pc += 2 },
	nil,
	nil,
	nil,
	func() { adc(immediate); pc += 2 },
	nil,
	nil,
	func() { jmp(indirectAddress) },
	func() { adc(absolute); pc += 3 },
	nil,
	nil,

	// 0x70 - 0x7F
	bvs,
	func() { adc(indirectIndexed); pc += 2 },
	nil,
	nil,
	nil,
	func() { adc(zeroPageX); pc += 2 },
	nil,
	nil,
	nil,
	func() { adc(absoluteY); pc += 3 },
	nil,
	nil,
	nil,
	func() { adc(absoluteX); pc += 3 },
	nil,
	nil,

	// 0x80 - 0x8F
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	dey,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0x90 - 0x9F
	bcc,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	// 0xA0 - 0xAF
	func() { ldy(immediate); pc += 2 },
	func() { lda(indexedIndirect); pc += 2 },
	func() { ldx(immediate); pc += 2 },
	nil,
	func() { ldy(zeroPage); pc += 2 },
	func() { lda(zeroPage); pc += 2 },
	func() { ldx(zeroPage); pc += 2 },
	nil,
	nil,
	func() { lda(immediate); pc += 2 },
	nil,
	nil,
	func() { ldy(absolute); pc += 3 },
	func() { lda(absolute); pc += 3 },
	func() { ldx(absolute); pc += 3 },
	nil,

	// 0xB0 - 0xBF
	bcs,
	func() { lda(indirectIndexed); pc += 2 },
	nil,
	nil,
	func() { ldy(zeroPageX); pc += 2 },
	func() { lda(zeroPageX); pc += 2 },
	func() { ldx(zeroPageY); pc += 2 },
	nil,
	clv,
	func() { lda(absoluteY); pc += 3 },
	nil,
	nil,
	func() { ldy(absoluteX); pc += 3 },
	func() { lda(absoluteX); pc += 3 },
	func() { ldx(absoluteY); pc += 3 },
	nil,

	// 0xC0 - 0xCF
	func() { cpy(immediate); pc += 2 },
	func() { cmp(indexedIndirect); pc += 2 },
	nil,
	nil,
	func() { cpy(zeroPage); pc += 2 },
	func() { cmp(zeroPage); pc += 2 },
	func() { dec(zeroPage); pc += 2 },
	nil,
	iny,
	func() { cmp(immediate); pc += 2 },
	dex,
	nil,
	func() { cpy(absolute); pc += 3 },
	func() { cmp(absolute); pc += 3 },
	func() { dec(absolute); pc += 3 },
	nil,

	// 0xD0 - 0xDF
	bne,
	func() { cmp(indirectIndexed); pc += 2 },
	nil,
	nil,
	nil,
	func() { cmp(zeroPageX); pc += 2 },
	func() { dec(zeroPageX); pc += 2 },
	nil,
	cld,
	func() { cmp(absoluteY); pc += 3 },
	nil,
	nil,
	nil,
	func() { cmp(absoluteX); pc += 3 },
	func() { dec(absoluteX); pc += 3 },
	nil,

	// 0xE0 - 0xEF
	func() { cpx(immediate); pc += 2 },
	nil,
	nil,
	nil,
	func() { cpx(zeroPage); pc += 2 },
	nil,
	func() { inc(zeroPage); pc += 2 },
	nil,
	inx,
	nil,
	nop,
	nil,
	func() { cpx(absolute); pc += 3 },
	nil,
	func() { inc(absolute); pc += 3 },
	nil,

	// 0xF0 - 0xFF
	beq,
	nil,
	nil,
	nil,
	nil,
	nil,
	func() { inc(zeroPageX); pc += 2 },
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	func() { inc(absoluteX); pc += 3 },
	nil}
