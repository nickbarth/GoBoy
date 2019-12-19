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

  // fmt.Printf("-- 0x%0x -- \n", cpu.sp)
  fmt.Printf("0x%0x\t", opcode)

  switch opcode {
  case 0x00:
    // nop
    fmt.Printf("0x%0x:\tnop\n", n,)
    cpu.t += 4
    cpu.pc += 1;

  case 0x05:
    // dec b
    cpu.t += 4
    cpu.reg.b -= 1

    cpu.reg.SetFlag('z', cpu.reg.b == 0)
    cpu.reg.SetFlag('n', true)
    cpu.reg.SetFlag('h', cpu.reg.b & 0x10 == 0x10)

    fmt.Printf("0x%0x:\tdec b\n", n)
    cpu.pc += 1

  case 0x06:
    // ld b nn
    cpu.t += 8
    cpu.reg.b = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("0x%0x:\tld b 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x08:
    // ld (nn) sp
    cpu.t += 8
    cpu.reg.c = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("0x%0x:\tld c 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x17:
    // rla
    cpu.t += 4
    cpu.reg.SetFlag('n', false)
    cpu.reg.SetFlag('h', false)
    c := cpu.reg.GetFlagVal('c')
    cpu.reg.SetFlag('c', cpu.reg.a & (1 << 7) >> 7 == 1)
    cpu.reg.a = cpu.reg.a << 1 | c
    fmt.Printf("0x%0x:\trla\n", n)
    cpu.pc += 1

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
    cpu.reg.Set16("de", uint16(cpu.mmu.ReadWord(cpu.pc + 1)))
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
    cpu.reg.SetFlag('n', false)
    cpu.reg.SetFlag('h', false)
    c := cpu.reg.GetFlagVal('c')
    cpu.reg.SetFlag('c', cpu.reg.a & 1 == 1)
    cpu.reg.a = cpu.reg.a >> 1 | (c << 7)
    fmt.Printf("0x%0x:\trra\n", n)
    cpu.pc += 1

  case 0x20:
    // jr nz nn
    fmt.Printf("0x%0x:\tjr nz 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
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
    fmt.Printf("0x%0x:\tld hl 0x%x\n", n, cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x22:
    // ldi (hl) hl -- gb sp
    cpu.t += 8
    cpu.mmu.Write(cpu.reg.Get16("hl"), cpu.reg.a)
    fmt.Printf("0x%0x:\tldi (hl) a\n", n)
    cpu.pc += 1

  case 0x23:
    // inc hl
    cpu.t += 8
    hl := cpu.reg.Get16("hl")
    cpu.reg.Set16("hl", hl+1)
    fmt.Printf("0x%0x:\thl inc\n", n)
    cpu.pc += 1

  case 0x31:
    // ld sp nnnn
    cpu.t += 12
    cpu.sp = uint16(cpu.mmu.ReadWord(cpu.pc + 1))
    fmt.Printf("0x%0x:\tld sp 0x%x\n", n, cpu.mmu.ReadWord(cpu.pc + 1))
    cpu.pc += 3

  case 0x32:
    // ld (hl-) a -- gb sp
    cpu.t += 8
    hl := cpu.reg.Get16("hl")
    cpu.mmu.Write(hl, cpu.reg.a)
    cpu.reg.Set16("hl", hl - 1)
    fmt.Printf("0x%0x:\tld (hl-) a\n", n)
    cpu.pc += 1

  case 0x3e:
    // ld a nn
    cpu.t += 8
    cpu.reg.a = cpu.mmu.Read(cpu.pc + 1)
    fmt.Printf("0x%0x:\tld a 0x%x\n", n, cpu.mmu.Read(cpu.pc + 1))
    cpu.pc += 2

  case 0x4f:
    // ld c a
    cpu.t += 4
    cpu.reg.c = cpu.reg.a
    fmt.Printf("0x%0x:\tld c a\n", n)
    cpu.pc += 1

  case 0x77:
    // ld (hl-) a
    cpu.t += 8
    hl := cpu.reg.Get16("hl")
    cpu.mmu.Write(hl, cpu.reg.a)
    cpu.reg.Set16("hl", hl - 1)
    fmt.Printf("0x%0x:\tld (hl-) a\n", n)
    cpu.pc += 1

  case 0x9c:
    // sbc a h
    cpu.t += 4
    val := cpu.reg.h + cpu.reg.GetFlagVal('c')

    cpu.reg.SetFlag('z', val == cpu.reg.a)
    cpu.reg.SetFlag('n', true)
    cpu.reg.SetFlag('c', val > cpu.reg.a)

    cpu.reg.a -=  val

    cpu.reg.SetFlag('h', cpu.reg.a & 0x10 == 0x10)
    fmt.Printf("0x%0x:\tsbc a h\n", n)
    cpu.pc += 1

  case 0xaf:
    // xor a
    cpu.t += 4
    cpu.reg.a = 0
    cpu.reg.f = 0x80
    fmt.Printf("0x%0x:\txor a\n", n)
    cpu.pc += 1

  case 0xc1:
    // pop bc
    cpu.t += 12
    cpu.reg.Set16("bc", cpu.sp)
    cpu.sp += 2
    fmt.Printf("0x%0x:\tpop bc\n", n)
    cpu.pc += 1

  case 0xcb:
    // bit instr
    instr := cpu.mmu.Read(cpu.pc + 1)

    switch instr {
    case 0x11:
      // rl c
      cpu.t += 8
      cpu.reg.SetFlag('n', false)
      cpu.reg.SetFlag('h', false)
      c := cpu.reg.GetFlagVal('c')
      cpu.reg.SetFlag('c', cpu.reg.c & (1 << 7) >> 7 == 1)
      cpu.reg.c = cpu.reg.c << 1 | c
      fmt.Printf("0x%0x:\trl c\n", n)
      cpu.pc += 2

    case 0x7c:
      // bit 7 h
      cpu.t += 8
      bit := cpu.reg.h & (1 << 6) >> 6
      cpu.reg.SetFlag('z', bit == 0)
      cpu.reg.SetFlag('n', false)
      cpu.reg.SetFlag('h', false)
      fmt.Printf("0x%0x:\tbit 7 h\n", n)
      cpu.pc += 2

      default:
        // not implemented
        panic(fmt.Sprintf("0x%0x:\tNI [0xcb%x]\n", n, instr))
    }

  case 0xc5:
    // push bc
    cpu.t += 16
    bc := cpu.reg.Get16("bc")
    cpu.mmu.WriteWord(cpu.sp - 2, bc)
    fmt.Printf("0x%0x:\tpush bc --\n", n)
    cpu.sp -= 2
    cpu.pc += 1

  case 0xc9:
    // ret
    cpu.t += 16
    fmt.Printf("0x%0x:\tret\n", n)
    cpu.pc = uint16(cpu.mmu.ReadWord(cpu.sp))
    cpu.sp += 2

  case 0xcd:
    // call nnnn
    cpu.t += 24
    cpu.mmu.WriteWord(cpu.sp - 2, cpu.pc + 3)
    nn := uint16(cpu.mmu.ReadWord(cpu.pc + 1))
    fmt.Printf("0x%0x:\tcall 0x%x ==\n", n, nn)
    cpu.sp -= 2
    cpu.pc = nn

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
