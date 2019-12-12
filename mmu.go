package main

import "io/ioutil"

type MMU struct {
  rom []byte
}

func NewMMU() *MMU {
  return &MMU{
    rom: []byte{},
  }
}

func (m *MMU) Load(file string) {
  dat, _ := ioutil.ReadFile(file)
  m.rom = dat
}

func (m *MMU) Read(addr uint16) byte {
  return m.rom[addr]
}

func (m *MMU) ReadWord(addr uint16) uint16 {
  return uint16(m.rom[addr] << 4 | m.rom[addr+1])
}

func (m *MMU) Write(addr int16, val byte) {
  m.rom[addr] = val
}

func (m *MMU) WriteWord(addr int16, val uint16) {
  m.rom[addr] = byte(val >> 4)
  m.rom[addr+1] = byte(val & 0b0000_1111)
}
