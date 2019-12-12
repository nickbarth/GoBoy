package main

import "fmt"

type Clock struct {
  m uint8
  t uint8
}

type CPU struct {
  reg *Registers
  mmu *MMU
  pc uint16 // program counter
  sp uint16 // stack pointer
  clock Clock
}

func NewCPU() *CPU {
  return &CPU{
    reg: NewRegisters(),
    mmu: NewMMU(),
    clock: Clock{ m: 0, t: 0},
    pc: 0,
    sp: 0,
  }
}

func (cpu *CPU) Step(n int) {
  opcode := cpu.mmu.Read(cpu.pc)
  cpu.pc += 1

  // fmt.Printf("Line %d\n", cpu.pc)

  switch opcode {
  case 0x00:
    // nop
    cpu.clock.m += 1
    cpu.clock.t += 4

  case 0xc3:
    // jp nn
    cpu.pc = cpu.mmu.ReadWord(cpu.pc)
    cpu.clock.t += 10

  case 0x31:
    // ld nn
    cpu.sp = cpu.mmu.ReadWord(cpu.pc)
    cpu.pc += 2
    cpu.clock.t += 10

  default:
    fmt.Printf("0x%x Not Found.\n", opcode)
  }
}
