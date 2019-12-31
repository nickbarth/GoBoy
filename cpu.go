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

// cpu instr
func (cpu *CPU) XOR(reg1 uint8, reg2 uint8) uint8 {
  val := reg1 ^ reg2
  if val == 0 {
   cpu.reg.SetFlag('z', true)
  }
  cpu.reg.SetFlag('n', false)
  cpu.reg.SetFlag('h', false)
  cpu.reg.SetFlag('c', false)
  return val
}

func (cpu *CPU) BIT(b uint8, reg uint8) uint8 {
  val := reg & (1 << b) >> b
  if val == 0 {
   cpu.reg.SetFlag('z', true)
  }
  cpu.reg.SetFlag('n', false)
  cpu.reg.SetFlag('h', true)
  return val
}

func (cpu *CPU) INC(reg uint8) uint8 {
  val := uint8((reg + 1) & 0xff)
  if val == 0 {
   cpu.reg.SetFlag('z', true)
  }
  cpu.reg.SetFlag('n', false)
  cpu.reg.SetFlag('h', (val & 0xF) == 0)
  return val
}

func (cpu *CPU) INC16(reg uint16) uint16 {
  val := uint16((reg + 1) & 0xffff)
  return val
}

func (cpu *CPU) DEC(reg uint8) uint8 {
  val := uint8((reg - 1) & 0xff)
  if val == 0 {
   cpu.reg.SetFlag('z', true)
  }
  cpu.reg.SetFlag('n', true)
  cpu.reg.SetFlag('h', (val & 0xF) == 0xf)
  return val
}

func (cpu *CPU) SBC(reg1 uint8, reg2 uint8) uint8 {
  c := cpu.reg.GetFlagVal('c')
  val := reg1 - reg2 - c
  if val == 0 {
   cpu.reg.SetFlag('z', true)
  }
  cpu.reg.SetFlag('n', true)
  cpu.reg.SetFlag('h', (reg1 & 0xF) - (reg2 & 0xF) - c < 0x0)
  cpu.reg.SetFlag('c', reg2 + c > reg1)
  return val & 0xFF
}

func (cpu *CPU) RL(reg uint8) uint8 {
  val := reg << 1 & 0xff
  if val == 0 {
   cpu.reg.SetFlag('z', true)
  }
  cpu.reg.SetFlag('n', false)
  cpu.reg.SetFlag('h', false)
  cpu.reg.SetFlag('c', (reg << 1) >> 8 == 1)
  return val
}

func (cpu *CPU) RR(reg uint8) uint8 {
  val := reg >> 1
  if val == 0 {
   cpu.reg.SetFlag('z', true)
  }
  cpu.reg.SetFlag('n', false)
  cpu.reg.SetFlag('h', false)
  cpu.reg.SetFlag('c', reg & 1 == 1)
  return val
}

func (cpu *CPU) POP() uint16 {
    val := cpu.mmu.ReadWord(cpu.sp)
    cpu.sp += 2
    return val
}

func (cpu *CPU) PUSH(val uint16) {
    cpu.sp -= 2
    cpu.mmu.WriteWord(cpu.sp, val)
}

func (cpu *CPU) Step(n uint16) {
  opcode := cpu.mmu.Read(cpu.pc)

  fmt.Printf("0x%x: \t(0x%0x)\t", n, opcode)

  switch opcode {
  case 0x00:
    // nop
    fmt.Printf("nop")
    cpu.t += 4
    cpu.pc += 1;

  case 0x05:
    // dec b
    cpu.t += 4
    cpu.reg.b = cpu.DEC(cpu.reg.b)
    fmt.Printf("dec b")
    cpu.pc += 1

  case 0x06:
    // ld b nn
    cpu.t += 8
    cpu.reg.b = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("ld b 0x%x", cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x08:
    // ld (nn) sp
    cpu.t += 8
    cpu.reg.c = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("ld c 0x%x", cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x17:
    // rla
    cpu.t += 4
    cpu.reg.a = cpu.RL(cpu.reg.a)
    fmt.Printf("rla")
    cpu.pc += 1

  case 0x0c:
    // inc c
    cpu.t += 4
    cpu.reg.c = cpu.INC(cpu.reg.c)
    fmt.Printf("inc c")
    cpu.pc += 1

  case 0x0e:
    // ld c nn
    cpu.t += 8
    cpu.reg.c = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("ld c 0x%x", cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x11:
    // ld de nnnn
    cpu.t += 12
    cpu.reg.Set16("de", uint16(cpu.mmu.ReadWord(cpu.pc + 1)))
    fmt.Printf("ld de 0x%x", cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x1a:
    // ld a,de
    cpu.t += 8
    cpu.reg.a = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("ld a de")
    cpu.pc += 1

  case 0x1f:
    // rra
    cpu.t += 4
    cpu.reg.a = cpu.RR(cpu.reg.a)
    fmt.Printf("rra")
    cpu.pc += 1

  case 0x20:
    // jr nz nn
    fmt.Printf("jr nz 0x%x", cpu.mmu.Read(cpu.pc + 1))
    if !cpu.reg.GetFlag('z') {
      cpu.pc += uint16(cpu.mmu.Read(cpu.pc + 1))
      cpu.t += 12
    } else {
      cpu.t += 8
      cpu.pc += 2
    }

  case 0x21:
    // ld hl nnnn
    cpu.t += 12
    cpu.reg.Set16("hl", uint16(cpu.mmu.ReadWord(cpu.pc + 1)))
    fmt.Printf("ld hl 0x%x", cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x22:
    // ldi (hl) hl -- gb sp
    cpu.t += 8
    cpu.mmu.Write(cpu.reg.Get16("hl"), cpu.reg.a)
    fmt.Printf("ldi (hl) a")
    cpu.pc += 1

  case 0x23:
    // inc hl
    cpu.t += 8
    hl := cpu.reg.Get16("hl")
    cpu.reg.Set16("hl", cpu.INC16(hl))
    fmt.Printf("hl inc")
    cpu.pc += 1

  case 0x31:
    // ld sp nnnn
    cpu.t += 12
    cpu.sp = uint16(cpu.mmu.ReadWord(cpu.pc + 1))
    fmt.Printf("ld sp 0x%x", cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x32:
    // ld (hl-) a -- gb sp
    cpu.t += 8
    hl := cpu.reg.Get16("hl")
    cpu.mmu.Write(hl, cpu.reg.a)
    cpu.reg.Set16("hl", hl - 1)
    fmt.Printf("ld (hl-) a")
    cpu.pc += 1

  case 0x3e:
    // ld a nn
    cpu.t += 8
    cpu.reg.a = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("ld a 0x%x", cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x4f:
    // ld c a
    cpu.t += 4
    cpu.reg.c = cpu.reg.a
    fmt.Printf("ld c a")
    cpu.pc += 1

  case 0x77:
    // ld (hl-) a
    cpu.t += 8
    hl := cpu.reg.Get16("hl")
    cpu.mmu.Write(hl, cpu.reg.a)
    cpu.reg.Set16("hl", hl - 1)
    fmt.Printf("ld (hl-) a")
    cpu.pc += 1

  case 0x9c:
    // sbc a h
    cpu.t += 4
    cpu.reg.a = cpu.SBC(cpu.reg.a, cpu.reg.h)
    fmt.Printf("sbc a h")
    cpu.pc += 1

  case 0xaf:
    // xor a
    cpu.t += 4
    cpu.reg.a = cpu.XOR(cpu.reg.a, cpu.reg.a)
    fmt.Printf("xor a")
    cpu.pc += 1

  case 0xc1:
    // pop bc
    cpu.t += 12
    cpu.reg.Set16("bc", cpu.POP())
    fmt.Printf("pop bc")
    cpu.pc += 1

  case 0xcb:
    // bit instr
    instr := cpu.mmu.Read(cpu.pc + 1)

    switch instr {
    case 0x11:
      // rl c
      cpu.t += 8
      cpu.reg.c = cpu.RL(cpu.reg.c)
      fmt.Printf("rl c")
      cpu.pc += 2

    case 0x7c:
      // bit 7 h
      cpu.t += 8
      cpu.BIT(7, cpu.reg.a)
      fmt.Printf("bit 7 h")
      cpu.pc += 2

      default:
        // not implemented
        panic(fmt.Sprintf("NI [0xcb%x]", instr))
    }

  case 0xc5:
    // push bc
    cpu.t += 16
    cpu.PUSH(cpu.reg.Get16("bc"))
    fmt.Printf("push bc")
    cpu.pc += 1

  case 0xc9:
    // ret
    cpu.t += 16
    fmt.Printf("ret")
    cpu.pc = uint16(cpu.mmu.ReadWord(cpu.sp))
    cpu.sp += 2

  case 0xcd:
    // call nnnn
    cpu.t += 24
    cpu.mmu.WriteWord(cpu.sp - 2, cpu.pc + 3)
    nn := uint16(cpu.mmu.ReadWord(cpu.pc + 1))
    fmt.Printf("call 0x%x", nn)
    cpu.sp -= 2
    cpu.pc = nn

  case 0xe0:
    // ld (ff00 + nn) a -- gb sp
    cpu.t += 8
    nn := uint16(cpu.mmu.Read(cpu.pc + 1))
    cpu.mmu.Write(0xff00 + nn, cpu.reg.a)
    fmt.Printf("ld (0xff00 + 0x%x) a", nn)
    cpu.pc += 2

  case 0xe2:
    // ld (FF00 + c) a -- gb sp
    cpu.t += 8
    cpu.mmu.Write(0xff00 + uint16(cpu.reg.c), cpu.reg.a)
    fmt.Printf("ld (0xff00 + c) a")
    cpu.pc += 1

  default:
    // not implemented
    panic(fmt.Sprintf("NI [0x%x]\n", opcode))
    // fmt.Printf("NI [0x%x]", opcode)
    cpu.pc += 1
  }

  fmt.Printf("\n")
}
