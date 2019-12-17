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
    fmt.Printf("0x%0x:\tnop\n", n,)
    cpu.t += 4
    cpu.pc += 1;

  case 0x08:
    // ld (nn) sp
    cpu.t += 8
    cpu.reg.c = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("0x%0x:\tld c 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x0c:
    // inc c
    cpu.t += 4
    cpu.reg.c += 1
    fmt.Printf("0x%0x:\tinc c\n", n)
    cpu.pc += 1

  case 0x0e:
    // ld c nn
    cpu.t += 8
    cpu.reg.c = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("0x%0x:\tld c 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x11:
    // ld de nnnn
    cpu.t += 12
    cpu.reg.set16("de", uint16(cpu.mmu.ReadWord(cpu.pc + 1)))
    fmt.Printf("0x%0x:\tld de 0x%x\n", n, cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x1a:
    // ld a,de
    cpu.t += 8
    cpu.reg.a = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("0x%0x:\tld a de\n", n)
    cpu.pc += 1

  case 0x1f:
    // rra
    cpu.t += 4
    cpu.reg.setFlag('n', false)
    cpu.reg.setFlag('h', false)
    cpu.reg.setFlag('c', cpu.reg.a & (1 << 7) >> 7 == 1)
    cpu.reg.a = cpu.reg.a &^ (1 << 7) << 1
    fmt.Printf("0x%0x:\trra\n", n)
    cpu.pc += 1

  case 0x20:
    // jr nz nn
    fmt.Printf("0x%0x:\tjr nz 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
    if !cpu.reg.getFlag('z') {
      cpu.pc += uint16(cpu.mmu.Read(cpu.pc + 1))
      cpu.t += 12
    } else {
      cpu.t += 8
      cpu.pc += 2
    }

  case 0x21:
    // ld hl nnnn
    cpu.t += 12
    cpu.reg.set16("hl", uint16(cpu.mmu.ReadWord(cpu.pc + 1)))
    fmt.Printf("0x%0x:\tld hl 0x%x\n", n, cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x31:
    // ld sp nnnn
    cpu.t += 12
    cpu.sp = uint16(cpu.mmu.ReadWord(cpu.pc + 1))
    fmt.Printf("0x%0x:\tld sp 0x%x\n", n, cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x32:
    // ld (hl-) a -- gb sp
    cpu.t += 8
    hl := cpu.reg.get16("hl")
    cpu.mmu.Write(hl, cpu.reg.a)
    cpu.reg.set16("hl", hl - 1)
    fmt.Printf("0x%0x:\tld (hl-) a\n", n)
    cpu.pc += 1

  case 0x3e:
    // ld a nn
    cpu.t += 8
    cpu.reg.a = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("0x%0x:\tld a 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x77:
    // ld (hl-) a
    cpu.t += 8
    hl := cpu.reg.get16("hl")
    cpu.mmu.Write(hl, cpu.reg.a)
    cpu.reg.set16("hl", hl - 1)
    fmt.Printf("0x%0x:\tld (hl-) a\n", n)
    cpu.pc += 1

  case 0xaf:
    // xor a
    cpu.t += 4
    cpu.reg.a = 0
    cpu.reg.f = 0x80
    fmt.Printf("0x%0x:\txor a\n", n)
    cpu.pc += 1

  case 0xcb:
    // bit instr
    instr := cpu.mmu.Read(cpu.pc + 1)

    switch instr {
    case 0x7c:
      cpu.t += 8
      bit := cpu.reg.h & (1 << 6) >> 6
      cpu.reg.setFlag('z', bit == 0)
      cpu.reg.setFlag('n', false)
      cpu.reg.setFlag('h', false)
      fmt.Printf("0x%0x:\tbit 7 h\n", n)
      cpu.pc += 2

      default:
        // not implemented
        panic(fmt.Sprintf("0x%0x:\tNI [0xcb%x]\n", n, instr))
    }

  case 0xe0:
    // ld (ff00 + nn) a -- gb sp
    cpu.t += 8
    nn := uint16(cpu.mmu.Read(cpu.pc + 1))
    cpu.mmu.Write(0xff00 + nn, cpu.reg.a)
    fmt.Printf("0x%0x:\tld (0xff00 + 0x%x) a\n", n, nn)
    cpu.pc += 2

  case 0xe2:
    // ld (FF00 + c) a -- gb sp
    cpu.t += 8
    cpu.mmu.Write(0xff00 + uint16(cpu.reg.c), cpu.reg.a)
    fmt.Printf("0x%0x:\tld (0xff00 + c) a\n", n)
    cpu.pc += 1

  default:
    // not implemented
    panic(fmt.Sprintf("0x%0x:\tNI [0x%x]\n", n, opcode))
    // fmt.Printf("0x%0x:\tNI [0x%x]\n", n, opcode)
    cpu.pc += 1
  }
}
