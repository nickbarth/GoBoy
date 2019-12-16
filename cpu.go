package main

import "fmt"

type CPU struct {
  reg *Registers
  mmu *MMU
  pc uint16 // program counter
  sp uint16 // stack pointer
  t uint8   // clock
}

func NewCPU() *CPU {
  return &CPU{
    reg: NewRegisters(),
    mmu: NewMMU(),
    pc: 0,
    sp: 0,
    t: 0,
  }
}

func (cpu *CPU) Step(n uint16) {
  opcode := cpu.mmu.Read(cpu.pc)

  switch opcode {
  case 0x00:
    // nop
    fmt.Printf("0x%0x: nop\n", n,)
    cpu.t += 4
    cpu.pc += 1;

  case 0x20:
    // jr nz nn
    fmt.Printf("0x%0x: jr nz 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
    if !cpu.reg.getFlag('z') {
      cpu.pc += uint16(cpu.mmu.Read(cpu.pc + 1))
      cpu.t += 12
    } else {
      cpu.t += 8
      cpu.pc += 2
    }

  case 0x21:
    // ld hl nnnn
    cpu.t += 10
    cpu.reg.set16("hl", uint16(cpu.mmu.ReadWord(cpu.pc + 1)))
    fmt.Printf("0x%0x: ld hl 0x%x\n", n, cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x31:
    // ld sp nnnn
    cpu.t += 10
    cpu.sp = uint16(cpu.mmu.ReadWord(cpu.pc + 1))
    fmt.Printf("0x%0x: ld sp 0x%x\n", n, cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x32:
    // ld nnnn a
    cpu.t += 13
    cpu.mmu.Write(uint16(cpu.mmu.ReadWord(cpu.pc + 1)), cpu.reg.a)
    fmt.Printf("0x%0x: ld 0x%x a\n", n, cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0xaf:
    // xor a
    cpu.t += 4
    cpu.reg.a = 0
    cpu.reg.f = 0x80
    fmt.Printf("0x%0x: xor a\n", n)
    cpu.pc += 1

  default:
    // not implemented
    fmt.Printf("Line %d: %x \n", cpu.pc, opcode)
    cpu.pc += 1
  }
}
