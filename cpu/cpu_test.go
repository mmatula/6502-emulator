package cpu

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
	"fmt"
)

type TestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *TestSuite) SetupTest() {
	pc = 0
	sp = 0xFF
	a = 0
	x = 0
	y = 0
	setPs(0)
	ticksToNext = 1
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCpuTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) runOp() {
	CodeToOp[Memory[pc]]()
}

func (suite *TestSuite) testAccumulator(expected uint8) {
	zero = expected != 0
	negative = expected & 128 == 0

	suite.runOp()

	suite.Equal(expected, a)
	suite.Equal(expected == 0, zero)
	suite.Equal(expected & 128 != 0, negative)
}

func (suite *TestSuite) TestAdcDecimalToSetCarryFlag() {
	// ADC #28
	Memory[0x0000] = 0x69
	Memory[0x0001] = 0x28

	decimalMode = true
	a = 0x79

	suite.runOp()

	suite.Equal(uint8(0x07), a)
	suite.True(carry)
	suite.False(zero)
	suite.True(overflow)
	suite.True(negative)
	suite.Equal(2, ticksToNext)
}

func (suite *TestSuite) TestAdcDecimalToSetZeroFlag() {
	// ADC #0
	Memory[0x0000] = 0x69
	Memory[0x0001] = 0x00

	decimalMode = true

	suite.testAccumulator(0)

	suite.False(carry)
	suite.False(overflow)
	suite.False(negative)
	suite.Equal(2, ticksToNext)
}

func (suite *TestSuite) TestLdaImmediate() {
	// LDA #F1
	Memory[0x0000] = 0xA9
	Memory[0x0001] = 0xF1

	suite.testAccumulator(0xF1)

	suite.Equal(2, ticksToNext)
}

func (suite *TestSuite) TestLdaAbsoluteX() {
	// LDA $1080, X
	Memory[0x0000] = 0xBD
	Memory[0x0001] = 0x80
	Memory[0x0002] = 0x10

	Memory[0x1084] = 0x76

	x = 4

	suite.testAccumulator(Memory[0x1084])
	suite.Equal(4, ticksToNext)
}

func (suite *TestSuite) TestLdaAbsoluteXPageCross() {
	// LDA $1080, X
	Memory[0x0000] = 0xBD
	Memory[0x0001] = 0x80
	Memory[0x0002] = 0x10

	Memory[0x1101] = 0x76

	x = 0x81

	suite.testAccumulator(Memory[0x1101])
	suite.Equal(5, ticksToNext)
}

func (suite *TestSuite) TestJmpIndirectPageCross() {
	// JMP ($02FF)
	Memory[0x0000] = 0x6C
	Memory[0x0001] = 0xFF
	Memory[0x0002] = 0x02

	Memory[0x0200] = 0x3C
	Memory[0x02FF] = 0xFF
	Memory[0x0300] = 0xFF

	suite.runOp()

	suite.Equal(uint16(0x3CFF), pc)
	suite.Equal(5, ticksToNext)
}

func (suite *TestSuite) TestSbcSimple() {
	// SBC #03
	Memory[0x0000] = 0xE9
	Memory[0x0001] = 0x03

	a = 0x05
	carry = true

	suite.testAccumulator(0x02)

	suite.True(carry)
	suite.False(overflow)
	suite.Equal(2, ticksToNext)
}

func (suite *TestSuite) TestSbcCarry() {
	// SBC #06
	Memory[0x0000] = 0xE9
	Memory[0x0001] = 0x06

	a = 0x05
	carry = true

	suite.testAccumulator(0xFF)

	suite.False(carry)
	suite.False(overflow)
	suite.Equal(2, ticksToNext)
}

func (suite *TestSuite) TestSbcOverflow() {
	// SBC #00
	Memory[0x0000] = 0xE9
	Memory[0x0001] = 0x00

	a = 0x80

	suite.testAccumulator(0x7F)

	suite.True(carry)
	suite.True(overflow)
	suite.Equal(2, ticksToNext)
}

func (suite *TestSuite) TestSbcDecimalCarry() {
	// SBC #27
	Memory[0x0000] = 0xE9
	Memory[0x0001] = 0x27

	a = 0x09
	carry = true
	decimalMode = true

	suite.testAccumulator(0x82)

	suite.False(carry)
	suite.False(overflow)
	suite.Equal(2, ticksToNext)
}

func (suite *TestSuite) TestSbcDecimal() {
	// SBC #09
	Memory[0x0000] = 0xE9
	Memory[0x0001] = 0x09

	a = 0x27
	carry = true
	decimalMode = true

	suite.testAccumulator(0x18)

	suite.True(carry)
	suite.False(overflow)
	suite.Equal(2, ticksToNext)
}

func (suite *TestSuite) TestRun() {
	// LDA #27 (2 ticks)
	Memory[0x0000] = 0xA9
	Memory[0x0001] = 0x27
	// SEC (2 ticks)
	Memory[0x0002] = 0x38
	// SED (2 ticks)
	Memory[0x0003] = 0xF8
	// SBC #09 (2 ticks)
	Memory[0x0004] = 0xE9
	Memory[0x0005] = 0x09
	// fill second page with the result
	// STA $0100, X (5 ticks)
	Memory[0x0006] = 0x9D
	Memory[0x0007] = 0x00
	Memory[0x0008] = 0x01
	// INX (2 ticks)
	Memory[0x0009] = 0xE8
	// BNE -4 (3/2 ticks)
	Memory[0x000A] = 0xD0
	Memory[0x000B] = 0xFC
	// JAM
	Memory[0x000C] = 0x02

	start := time.Now()
	run(0x0000)
	total := time.Now().Sub(start)

	suite.True(total.Nanoseconds() > 2567 * int64(tickDuration))
	for i := 0x0100; i < 0x0200; i++ {
		suite.Equal(uint8(0x18), Memory[i], fmt.Sprintf("%#x", i))
	}

	suite.True(carry)
	suite.False(overflow)
}