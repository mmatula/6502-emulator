package main

import "github.com/6502-emulator/cpu"

func main() {
	// LDA $23
	cpu.Memory[0x0000] = 0xA5
	cpu.Memory[0x0001] = 0x00
	cpu.Memory[0x0002] = 0x00
	cpu.Memory[0x0003] = 0x00
	cpu.Memory[0x0004] = 0x00
	cpu.Memory[0x0005] = 0x00

	//...
}
