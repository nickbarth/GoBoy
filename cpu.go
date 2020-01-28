package main

import "fmt"

type CPU struct {
	reg *Registers
	mmu *MMU
	pc  uint16 // program counter
	sp  uint16 // stack pointer
	t   uint8  // clock
}

func NewCPU() *CPU {
	return &CPU{
		reg: NewRegisters(),
		mmu: NewMMU(),
		pc:  0,
		sp:  0,
		t:   0,
	}
}

func (cpu *CPU) XOR(reg1 uint8, reg2 uint8) uint8 {
	val := reg1 ^ reg2
	cpu.reg.SetFlag('z', val == 0)
	cpu.reg.SetFlag('n', false)
	cpu.reg.SetFlag('h', false)
	cpu.reg.SetFlag('c', false)
	return val
}

func (cpu *CPU) BIT(b uint8, reg uint8) uint8 {
	val := reg & (1 << b) >> b
	cpu.reg.SetFlag('z', val == 0)
	cpu.reg.SetFlag('n', false)
	cpu.reg.SetFlag('h', true)
	return val
}

func (cpu *CPU) INC(reg byte) {
	val := uint8((cpu.reg.Get(reg) + 1) & 0xff)
	cpu.reg.SetFlag('z', val == 0)
	cpu.reg.SetFlag('n', false)
	cpu.reg.SetFlag('h', (val&0xF) == 0)
	cpu.reg.Set(reg, val)
}

func (cpu *CPU) INC16(reg string) {
	val := cpu.reg.Get16(reg)
	cpu.reg.Set16(reg, uint16((val+1)&0xffff))
}

func (cpu *CPU) DEC(reg byte) {
	val := uint8((cpu.reg.Get(reg) - 1) & 0xff)
	cpu.reg.SetFlag('z', val == 0)
	cpu.reg.SetFlag('n', true)
	cpu.reg.SetFlag('h', (val&0xF) == 0xf)
	cpu.reg.Set(reg, val)
}

func (cpu *CPU) LD_REG(reg1 byte, reg2 byte) {
	cpu.reg.Set(reg1, cpu.reg.Get(reg2))
}

func (cpu *CPU) LD_VAL(reg1 byte, val uint8) {
	cpu.reg.Set(reg1, val)
}

func (cpu *CPU) DEC16(reg string) {
	val := cpu.reg.Get16(reg)
	cpu.reg.Set16(reg, val-1)
}

func (cpu *CPU) SUB(reg uint8) {
	val := cpu.reg.a - reg
	cpu.reg.SetFlag('z', val == 0)
	cpu.reg.SetFlag('n', true)
	if (reg & 0xf) > (cpu.reg.a & 0xf) {
		cpu.reg.SetFlag('h', true)
	}
	cpu.reg.SetFlag('n', true)
	if reg > cpu.reg.a {
		cpu.reg.SetFlag('c', true)
	}
	cpu.reg.a = val
}

func (cpu *CPU) SBC(reg1 uint8, reg2 uint8) uint8 {
	c := cpu.reg.GetFlagVal('c')
	val := reg1 - reg2 - c
	cpu.reg.SetFlag('z', val == 0)
	cpu.reg.SetFlag('n', true)
	cpu.reg.SetFlag('h', (reg1&0xF)-(reg2&0xF)-c < 0x0)
	cpu.reg.SetFlag('c', reg2+c > reg1)
	return val & 0xFF
}

func (cpu *CPU) RL(reg uint8) uint8 {
	val := reg << 1 & 0xff
	cpu.reg.SetFlag('z', val == 0)
	cpu.reg.SetFlag('n', false)
	cpu.reg.SetFlag('h', false)
	cpu.reg.SetFlag('c', (reg<<1)>>8 == 1)
	return val
}

func (cpu *CPU) RR(reg uint8) uint8 {
	val := reg >> 1
	cpu.reg.SetFlag('z', val == 0)
	cpu.reg.SetFlag('n', false)
	cpu.reg.SetFlag('h', false)
	cpu.reg.SetFlag('c', reg&1 == 1)
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

	if opcode != 0xcb {
		fmt.Printf("0x%x: \t(0x%0x)\t\t", n, opcode)
	} else {
		fmt.Printf("0x%x: \t(0x%0x", n, opcode)
	}

	switch opcode {
	case 0x00:
		// nop
		fmt.Printf("nop")
		cpu.t += 4
		cpu.pc += 1

	case 0x04:
		// inc b
		cpu.t += 4
		cpu.INC('b')
		fmt.Printf("inc b")
		cpu.pc += 1

	case 0x05:
		// dec b
		cpu.t += 4
		cpu.DEC('b')
		fmt.Printf("dec b")
		cpu.pc += 1

	case 0x06:
		// ld b nn
		cpu.t += 8
		cpu.reg.b = cpu.mmu.Read(cpu.pc + 1)
		fmt.Printf("ld b 0x%x", cpu.mmu.Read(cpu.pc+1))
		cpu.pc += 2

	case 0x08:
		// ld (nn) sp
		cpu.t += 8
		cpu.reg.c = cpu.mmu.Read(cpu.pc + 1)
		fmt.Printf("ld c 0x%x", cpu.mmu.Read(cpu.pc+1))
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
		cpu.INC('c')
		fmt.Printf("inc c")
		cpu.pc += 1

	case 0x0d:
		// dec c
		cpu.t += 4
		cpu.DEC('c')
		fmt.Printf("dec c")
		cpu.pc += 1

	case 0x0e:
		// ld c nn
		cpu.t += 8
		cpu.reg.c = cpu.mmu.Read(cpu.pc + 1)
		fmt.Printf("ld c 0x%x", cpu.mmu.Read(cpu.pc+1))
		cpu.pc += 2

	case 0x11:
		// ld de nnnn
		cpu.t += 12
		cpu.reg.Set16("de", uint16(cpu.mmu.ReadWord(cpu.pc+1)))
		fmt.Printf("ld de 0x%x", cpu.mmu.ReadWord(cpu.pc+1))
		cpu.pc += 3

	case 0x13:
		// inc de
		cpu.t += 8
		cpu.INC16("de")
		fmt.Printf("inc de")
		cpu.pc += 1

	case 0x15:
		// dec d
		cpu.t += 4
		cpu.DEC('d')
		fmt.Printf("dec d")
		cpu.pc += 1

	case 0x16:
		// ld d nn
		cpu.t += 8
		nn := cpu.mmu.Read(cpu.pc + 1)
		cpu.LD_VAL('d', nn)
		fmt.Printf("ld d 0x%x", nn)
		cpu.pc += 2

	case 0x18:
		// jr nn
		cpu.t += 12
		nn := byte(cpu.mmu.Read(cpu.pc + 1))
		fmt.Printf("jr 0x%x", nn)
		cpu.pc += uint16(int8(nn)) + 2

	case 0x1a:
		// ld a (de)
		cpu.t += 8
		cpu.reg.a = cpu.mmu.Read(cpu.reg.Get16("de"))
		fmt.Printf("ld a (de)")
		cpu.pc += 1

	case 0x1d:
		// dec e
		cpu.t += 4
		cpu.DEC('e')
		fmt.Printf("dec e")
		cpu.pc += 1

	case 0x1e:
		// ld e nn
		cpu.t += 8
		cpu.reg.e = cpu.mmu.Read(cpu.pc + 1)
		fmt.Printf("ld e 0x%x", cpu.mmu.Read(cpu.pc+1))
		cpu.pc += 2

	case 0x1f:
		// rra
		cpu.t += 4
		cpu.reg.a = cpu.RR(cpu.reg.a)
		fmt.Printf("rra")
		cpu.pc += 1

	case 0x20:
		// jr nz nn
		nn := byte(cpu.mmu.Read(cpu.pc + 1))
		fmt.Printf("jr nz 0x%x", nn)
		if !cpu.reg.GetFlag('z') {
			cpu.pc += uint16(int8(nn)) + 2
			cpu.t += 12
		} else {
			cpu.t += 8
			cpu.pc += 2
		}

	case 0x21:
		// ld hl nnnn
		cpu.t += 12
		cpu.reg.Set16("hl", uint16(cpu.mmu.ReadWord(cpu.pc+1)))
		fmt.Printf("ld hl 0x%x", cpu.mmu.ReadWord(cpu.pc+1))
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
		cpu.INC16("hl")
		fmt.Printf("hl inc")
		cpu.pc += 1

	case 0x24:
		// inc h
		cpu.t += 4
		cpu.INC('h')
		fmt.Printf("inc h")
		cpu.pc += 1

	case 0x28:
		// jr z nn
		fmt.Printf("jr z nn")
		if cpu.reg.GetFlag('z') {
			nn := cpu.mmu.Read(cpu.pc + 1)
			cpu.pc = cpu.pc + 2 + uint16(nn)
			cpu.t += 12
		} else {
			cpu.t += 8
			cpu.pc += 2
		}

	case 0x2e:
		// ld l nn
		cpu.t += 8
		nn := uint8(cpu.mmu.Read(cpu.pc + 1))
		fmt.Printf("ld l 0x%x", nn)
		cpu.reg.l = nn
		cpu.pc += 2

	case 0x31:
		// ld sp nnnn
		cpu.t += 12
		cpu.sp = uint16(cpu.mmu.ReadWord(cpu.pc + 1))
		fmt.Printf("ld sp 0x%x", cpu.mmu.ReadWord(cpu.pc+1))
		cpu.pc += 3

	case 0x32:
		// ldd (hl) a -- gb sp
		cpu.t += 8
		hl := cpu.reg.Get16("hl")
		fmt.Printf("\n -- 0x%x -- \n", cpu.mmu.Read(hl))

		cpu.mmu.Write(hl, cpu.reg.a)
		cpu.mmu.Write(hl, 0x1)

		fmt.Printf("\n -- 0x%x -- \n", cpu.mmu.Read(hl))
		cpu.DEC16("hl")
		fmt.Printf("ldd (hl) a")
		cpu.pc += 1

	case 0x34:
		// inc (hl)
		cpu.t += 12
		cpu.INC16("hl")
		fmt.Printf("inc (hl)")
		cpu.pc += 1

	case 0x3e:
		// ld a nn
		cpu.t += 8
		cpu.reg.a = cpu.mmu.Read(cpu.pc + 1)
		fmt.Printf("ld a 0x%x", cpu.mmu.Read(cpu.pc+1))
		cpu.pc += 2

	case 0x3d:
		// dec a
		cpu.t += 4
		cpu.DEC('a')
		fmt.Printf("dec a")
		cpu.pc += 1

	case 0x4f:
		// ld c a
		cpu.t += 4
		cpu.LD_REG('c', 'a')
		fmt.Printf("ld c a")
		cpu.pc += 1

	case 0x57:
		// ld d a
		cpu.t += 4
		cpu.LD_REG('d', 'a')
		fmt.Printf("ld d a")
		cpu.pc += 1

	case 0x62:
		// ld h c
		cpu.t += 4
		cpu.LD_REG('h', 'c')
		fmt.Printf("ld h c")
		cpu.pc += 1

	case 0x67:
		// ld h a
		cpu.t += 4
		cpu.LD_REG('h', 'a')
		fmt.Printf("ld h a")
		cpu.pc += 1

	case 0x77:
		// ld (hl) a
		cpu.t += 8
		hl := cpu.reg.Get16("hl")
		cpu.mmu.Write(hl, cpu.reg.a)
		fmt.Printf("ld (hl) a")
		cpu.pc += 1

	case 0x7b:
		// ld a e
		cpu.t += 4
		cpu.LD_REG('a', 'e')
		fmt.Printf("ld a e")
		cpu.pc += 1

	case 0x7c:
		// ld a h
		cpu.t += 4
		cpu.LD_REG('a', 'h')
		fmt.Printf("ld a h")
		cpu.pc += 1

	case 0x90:
		// sub b
		cpu.t += 4
		cpu.SUB(cpu.reg.b)
		fmt.Printf("sub b")
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
		fmt.Printf("%x)\t", instr)

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
			cpu.BIT(7, cpu.reg.h)
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
		cpu.mmu.WriteWord(cpu.sp-2, cpu.pc+3)
		nn := uint16(cpu.mmu.ReadWord(cpu.pc + 1))
		fmt.Printf("call 0x%x", nn)
		cpu.sp -= 2
		cpu.pc = nn

	case 0xe0:
		// ld (ff00 + nn) a -- gb sp
		cpu.t += 8
		nn := uint16(cpu.mmu.Read(cpu.pc + 1))
		cpu.mmu.Write(0xff00|nn, cpu.reg.a)
		fmt.Printf("ld (0xff00 + 0x%x) a", nn)
		cpu.pc += 2

	case 0xe2:
		// ld (FF00 + c) a -- gb sp
		cpu.t += 8
		cpu.mmu.Write(0xff00+uint16(cpu.reg.c), cpu.reg.a)
		fmt.Printf("ld (0xff00 + c) a")
		cpu.pc += 1

	case 0xea:
		// ld (nnnn) a
		cpu.t += 16
		nnnn := cpu.mmu.ReadWord(cpu.pc + 1)
		cpu.mmu.Write(nnnn, cpu.reg.a)
		fmt.Printf("ld (nnnn) a")
		cpu.pc += 3

	case 0xf0:
		// ld (ff00 + nn) a -- gb sp
		cpu.t += 12
		nn := uint16(cpu.mmu.Read(cpu.pc + 1))
		cpu.reg.a = cpu.mmu.Read(0xff00 + nn)
		fmt.Printf("ld a (0xff00 + 0x%x)", nn)
		cpu.pc += 2

	case 0xfe:
		// cp a n
		cpu.t += 8
		cpu.mmu.Write(0xff00+uint16(cpu.reg.c), cpu.reg.a)
		fmt.Printf("cp a n")
		cpu.pc += 1

	default:
		// not implemented
		panic(fmt.Sprintf("NI [0x%x]\n", opcode))
		// fmt.Printf("NI [0x%x]", opcode)
		cpu.pc += 1
	}

	fmt.Printf("\n")
}

func (cpu *CPU) Run(n int64) {
	cpu.Step(cpu.pc)
}
