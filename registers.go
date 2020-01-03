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

func (r *Registers) Set16(reg16 string, val uint16) {
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

func (r Registers) Get16(reg16 string) uint16 {
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

func (r Registers) Get(reg byte) uint8 {
  switch reg {
    case 'a': return r.a
    case 'b': return r.b
    case 'c': return r.c
    case 'd': return r.d
    case 'e': return r.e
    case 'f': return r.f
    case 'h': return r.h
    case 'l': return r.l
  }
  return 0
}

func (r Registers) Set(reg byte, val uint8) {
  switch reg {
    case 'a': r.a = val
    case 'b': r.b = val
    case 'c': r.c = val
    case 'd': r.d = val
    case 'e': r.e = val
    case 'f': r.f = val
    case 'h': r.h = val
    case 'l': r.l = val
  }
}

// flag register
// z n h c 0 0 0 0
// z - zero flag
// n - substract flag
// h - half carry flag
// c - carry flag

func (r *Registers) SetFlag(flag byte, val bool) {
  switch flag {
  case 'z':
    if val { r.f = r.f |  (1 << 7)
    } else { r.f = r.f &^ (1 << 7) }
  case 'n':
    if val { r.f = r.f |  (1 << 6)
    } else { r.f = r.f &^ (1 << 6) }
  case 'h':
    if val { r.f = r.f |  (1 << 5)
    } else { r.f = r.f &^ (1 << 5) }
  case 'c':
    if val { r.f = r.f |  (1 << 4)
    } else { r.f = r.f &^ (1 << 4) }
  }
}

func (r *Registers) GetFlag(flag byte) bool {
  switch flag {
  case 'z':
    return r.f >> 7 == 1
  case 'n':
    return r.f & (1 << 6) >> 6 == 1
  case 'h':
    return r.f & (1 << 5) >> 5 == 1
  case 'c':
    return r.f & (1 << 4) >> 4 == 1
  }
  return false
}

func (r *Registers) GetFlagVal(flag byte) uint8 {
  switch flag {
  case 'z':
    return r.f >> 7
  case 'n':
    return r.f & (1 << 6) >> 6
  case 'h':
    return r.f & (1 << 5) >> 5
  case 'c':
    return r.f & (1 << 4) >> 4
  }
  return 0
}
