package main

import "io/ioutil"
import _"fmt"

type MMU struct {
  bios bool
  rom []byte
  ram []byte
}

func NewMMU() *MMU {
  return &MMU{
    bios: true,
    rom: []byte{},
    ram: []byte{},
  }
}

func (m *MMU) Load(file string) {
  dat, _ := ioutil.ReadFile(file)
  m.ram = append(dat, make([]byte, 0xffff)...) // hack
}

func (m *MMU) Read(addr uint16) byte {
  switch addr & 0xF000 {
  // bios // rom0
  case 0x0000:
    if m.bios {
      // fmt.Printf("BIOS -- ")
    } else {
      // fmt.Printf("ROM0 -- ")
    }
  // rom0
  case 0x1000:
  case 0x2000:
  case 0x3000:
    // fmt.Printf("ROM0 -- ")
  // rom1
  case 0x4000:
  case 0x5000:
  case 0x6000:
  case 0x7000:
    // fmt.Printf("ROM1 -- ")
  // vram
  case 0x8000:
  case 0x9000:
    // fmt.Printf("VRAM -- ")
  // eram
  case 0xa000:
  case 0xb000:
    // fmt.Printf("ERAM -- ")
  // wram
  case 0xc000:
  case 0xd000:
    // fmt.Printf("WRAM -- ")
  // wram shadow
  case 0xe000:
    // fmt.Printf("WRAM SD -- ")
  // shadow / io / zp ram
  case 0xf000:
    // fmt.Printf("IO -- ")
  default:
    // fmt.Printf("IO -- ")
  }
  return m.ram[addr]
}

func (m *MMU) ReadWord(addr uint16) uint16 {
  return uint16(m.Read(addr)) | uint16(m.Read(addr+1)) << 8
}

func (m *MMU) Write(addr uint16, val byte) {
  m.ram[addr] = val
}

func (m *MMU) WriteWord(addr uint16, val uint16) {
  m.ram[addr] = byte(val >> 4)
  m.ram[addr+1] = byte(val & 0b0000_1111)
}
