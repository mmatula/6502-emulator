package cpu

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type CpuTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *CpuTestSuite) SetupTest() {
	pc = 0
	sp = 0xFF
	a = 0
	x = 0
	y = 0
	setPs(0);
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCpuTestSuite(t *testing.T) {
	suite.Run(t, new(CpuTestSuite))
}

func (suite *CpuTestSuite) runOp() {
	CodeToOp[Memory[pc]]()
}

func (suite *CpuTestSuite) testAccumulator(expected uint8) {
	zero = expected != 0
	negative = expected & 128 == 0

	suite.runOp()

	suite.Equal(expected, a)
	suite.Equal(expected == 0, zero)
	suite.Equal(expected & 128 != 0, negative)
}

func (suite *CpuTestSuite) TestAdcDecimalToSetCarryFlag() {
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
}

func (suite *CpuTestSuite) TestAdcDecimalToSetZeroFlag() {
	// ADC #0
	Memory[0x0000] = 0x69
	Memory[0x0001] = 0x00

	decimalMode = true

	suite.testAccumulator(0)

	suite.False(carry)
	suite.False(overflow)
	suite.False(negative)
}

func (suite *CpuTestSuite) TestLdaImmediate() {
	// LDA $1080, X
	Memory[0x0000] = 0xA9
	Memory[0x0001] = 0xF1

	suite.testAccumulator(0xF1)
}

func (suite *CpuTestSuite) TestLdaAbsoluteX() {
	// LDA $1080, X
	Memory[0x0000] = 0xBD
	Memory[0x0001] = 0x80
	Memory[0x0002] = 0x10

	Memory[0x1084] = 0x76

	x = 4

	suite.testAccumulator(Memory[0x1084])
}

func (suite *CpuTestSuite) TestJmpIndirectPageCross() {
	// JMP ($02FF)
	Memory[0x0000] = 0x6C
	Memory[0x0001] = 0xFF
	Memory[0x0002] = 0x02

	Memory[0x0200] = 0x3C
	Memory[0x02FF] = 0xFF
	Memory[0x0300] = 0xFF

	suite.runOp()

	suite.Equal(uint16(0x3CFF), pc)
}