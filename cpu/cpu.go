package cpu

import "time"

// registers
var pc uint16 = 0		// program counter
var sp uint8 = 0xFF		// stack pointer
var a uint8				// accumulator
var x, y uint8			// index registers
var carry, zero, interruptDisable, decimalMode, breakCommand, overflow, negative bool	// processor status (flags)

const freq int64 = 985000
const tickDuration time.Duration = time.Duration(1000000000 / freq)
var ticksToNext = 0

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
	func() { ora(indexedIndirect) },
	nil,
	nil,
	nil,
	func() { ora(zeroPage) },
	func() { asl(zeroPage) },
	nil,
	php,
	func() { ora(immediate) },
	func() { asl(accumulator) },
	nil,
	nil,
	func() { ora(absolute) },
	func() { asl(absolute) },
	nil,

	// 0x10 - 0x1F
	bpl,
	func() { ora(indirectIndexed) },
	nil,
	nil,
	nil,
	func() { ora(zeroPageX) },
	func() { asl(zeroPageX) },
	nil,
	clc,
	func() { ora(absoluteY) },
	nil,
	nil,
	nil,
	func() { ora(absoluteX) },
	func() { asl(absoluteX) },
	nil,

	// 0x20 - 0x2F
	jsr,
	func() { and(indexedIndirect) },
	nil,
	nil,
	func() { bit(zeroPage) },
	func() { and(zeroPage) },
	func() { rol(zeroPage) },
	nil,
	plp,
	func() { and(immediate) },
	func() { rol(accumulator) },
	nil,
	func() { bit(absolute) },
	func() { and(absolute) },
	func() { rol(absolute) },
	nil,

	// 0x30 - 0x3F
	bmi,
	func() { and(indirectIndexed) },
	nil,
	nil,
	nil,
	func() { and(zeroPageX) },
	func() { rol(zeroPageX) },
	nil,
	sec,
	func() { and(absoluteY) },
	nil,
	nil,
	nil,
	func() { and(absoluteX) },
	func() { rol(absoluteX) },
	nil,

	// 0x40 - 0x4F
	rti,
	func() { eor(indexedIndirect) },
	nil,
	nil,
	nil,
	func() { eor(zeroPage) },
	func() { lsr(zeroPage) },
	nil,
	pha,
	func() { eor(immediate) },
	func() { lsr(accumulator) },
	nil,
	func() { jmp(absoluteAddress) },
	func() { eor(absolute) },
	func() { lsr(absolute) },
	nil,

	// 0x50 - 0x5F
	bvc,
	func() { eor(indirectIndexed) },
	nil,
	nil,
	nil,
	func() { eor(zeroPageX) },
	func() { lsr(zeroPageX) },
	nil,
	cli,
	func() { eor(absoluteY) },
	nil,
	nil,
	nil,
	func() { eor(absoluteX) },
	func() { lsr(absoluteX) },
	nil,

	// 0x60 - 0x6F
	rts,
	func() { adc(indexedIndirect) },
	nil,
	nil,
	nil,
	func() { adc(zeroPage) },
	func() { ror(zeroPage) },
	nil,
	pla,
	func() { adc(immediate) },
	func() { ror(accumulator) },
	nil,
	func() { jmp(indirectAddress) },
	func() { adc(absolute) },
	func() { ror(absolute) },
	nil,

	// 0x70 - 0x7F
	bvs,
	func() { adc(indirectIndexed) },
	nil,
	nil,
	nil,
	func() { adc(zeroPageX) },
	func() { ror(zeroPageX) },
	nil,
	sei,
	func() { adc(absoluteY) },
	nil,
	nil,
	nil,
	func() { adc(absoluteX) },
	func() { ror(absoluteX) },
	nil,

	// 0x80 - 0x8F
	nil,
	func() { sta(indexedIndirect) },
	nil,
	nil,
	func() { sty(zeroPage) },
	func() { sta(zeroPage) },
	func() { stx(zeroPage) },
	nil,
	dey,
	nil,
	txa,
	nil,
	func() { sty(absolute) },
	func() { sta(absolute) },
	func() { stx(absolute) },
	nil,

	// 0x90 - 0x9F
	bcc,
	func() { sta(indirectIndexed) },
	nil,
	nil,
	func() { sty(zeroPageX) },
	func() { sta(zeroPageX) },
	func() { stx(zeroPageY) },
	nil,
	tya,
	func() { sta(absoluteY) },
	txs,
	nil,
	nil,
	func() { sta(absoluteX) },
	nil,
	nil,

	// 0xA0 - 0xAF
	func() { ldy(immediate) },
	func() { lda(indexedIndirect) },
	func() { ldx(immediate) },
	nil,
	func() { ldy(zeroPage) },
	func() { lda(zeroPage) },
	func() { ldx(zeroPage) },
	nil,
	tay,
	func() { lda(immediate) },
	tax,
	nil,
	func() { ldy(absolute) },
	func() { lda(absolute) },
	func() { ldx(absolute) },
	nil,

	// 0xB0 - 0xBF
	bcs,
	func() { lda(indirectIndexed) },
	nil,
	nil,
	func() { ldy(zeroPageX) },
	func() { lda(zeroPageX) },
	func() { ldx(zeroPageY) },
	nil,
	clv,
	func() { lda(absoluteY) },
	tsx,
	nil,
	func() { ldy(absoluteX) },
	func() { lda(absoluteX) },
	func() { ldx(absoluteY) },
	nil,

	// 0xC0 - 0xCF
	func() { cpy(immediate) },
	func() { cmp(indexedIndirect) },
	nil,
	nil,
	func() { cpy(zeroPage) },
	func() { cmp(zeroPage) },
	func() { dec(zeroPage) },
	nil,
	iny,
	func() { cmp(immediate) },
	dex,
	nil,
	func() { cpy(absolute) },
	func() { cmp(absolute) },
	func() { dec(absolute) },
	nil,

	// 0xD0 - 0xDF
	bne,
	func() { cmp(indirectIndexed) },
	nil,
	nil,
	nil,
	func() { cmp(zeroPageX) },
	func() { dec(zeroPageX) },
	nil,
	cld,
	func() { cmp(absoluteY) },
	nil,
	nil,
	nil,
	func() { cmp(absoluteX) },
	func() { dec(absoluteX) },
	nil,

	// 0xE0 - 0xEF
	func() { cpx(immediate) },
	func() { sbc(indexedIndirect) },
	nil,
	nil,
	func() { cpx(zeroPage) },
	func() { sbc(zeroPage) },
	func() { inc(zeroPage) },
	nil,
	inx,
	func() { sbc(immediate) },
	nop,
	nil,
	func() { cpx(absolute) },
	func() { sbc(absolute) },
	func() { inc(absolute) },
	nil,

	// 0xF0 - 0xFF
	beq,
	func() { sbc(indirectIndexed) },
	nil,
	nil,
	nil,
	func() { sbc(zeroPageX) },
	func() { inc(zeroPageX) },
	nil,
	sed,
	func() { sbc(absoluteY) },
	nil,
	nil,
	nil,
	func() { sbc(absoluteX) },
	func() { inc(absoluteX) },
	nil}

func run(address uint16) {
	nextTick := time.Now()
	ticker := time.NewTicker(tickDuration)
	pc = address
	for CodeToOp[Memory[pc]] != nil {
		ticksToNext = 1
		CodeToOp[Memory[pc]]()
		nextTick = nextTick.Add(time.Duration(ticksToNext) * tickDuration)
		for time.Now().Before(nextTick) {
			<- ticker.C
		}
	}
}