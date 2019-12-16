package main

type Registers struct {
  a uint8
  b uint8
  c uint8
  d uint8
  e uint8
  f uint8
  h uint8
  l uint8
}

func NewRegisters() *Registers {
  return &Registers{
    a: 0,
    b: 0,
    c: 0,
    d: 0,
    e: 0,
    f: 0,
    h: 0,
    l: 0,
  }
}

func (r *Registers) set16(reg16 string, val uint16) {
  switch reg16 {
    case "af":
      r.a = uint8(val >> 8)
      r.f = uint8(val & 240) // lf always 0
    case "bc":
      r.b = uint8(val >> 8)
      r.c = uint8(val & 255)
    case "de":
      r.d = uint8(val >> 8)
      r.e = uint8(val & 255)
    case "hl":
      r.h = uint8(val >> 8)
      r.l = uint8(val & 255)
  }
}

func (r Registers) get16(reg16 string) uint16 {
  switch reg16 {
  case "af":
    return uint16(r.a) << 8 | uint16(r.f)
  case "bc":
    return uint16(r.b) << 8 | uint16(r.c)
  case "de":
    return uint16(r.d) << 8 | uint16(r.e)
  case "hl":
    return uint16(r.h) << 8 | uint16(r.l)
  }
  return 0;
}

// flag register
// z - zero flag
// n - substract flag
// h - half carry flag
// c - carry flag

func (r *Registers) setFlag(flag byte, val bool) {
  switch flag {
  case 'z':
    if val { r.f = r.f | 0b1000_0000
    } else { r.f = r.f & 0b0111_0000 }
  case 'n':
    if val { r.f = r.f | 0b0100_0000
    } else { r.f = r.f & 0b1011_0000 }
  case 'h':
    if val { r.f = r.f | 0b0010_0000
    } else { r.f = r.f & 0b1101_0000 }
  case 'c':
    if val { r.f = r.f | 0b0001_0000
    } else { r.f = r.f & 0b1110_0000 }
  }
}

func (r *Registers) getFlag(flag byte) bool {
  switch flag {
  case 'z':
    return r.f >> 7 == 1
  case 'n':
    return r.f & 0b0100_0000 >> 6 == 1
  case 'h':
    return r.f & 0b0100_0000 >> 5 == 1
  case 'c':
    return r.f & 0b0100_0000 >> 4 == 1
  }
  return false
}
