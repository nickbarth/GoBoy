package main

import "io/ioutil"

type MMU struct {
  ram []byte
}

func NewMMU() *MMU {
  return &MMU{
    ram: []byte{},
  }
}

func (m *MMU) Load(file string) {
  dat, _ := ioutil.ReadFile(file)
  m.ram = append(dat, make([]byte, 0xffff)...) // fix
}

func (m *MMU) Read(addr uint16) byte {
  return m.ram[addr]
}

func (m *MMU) ReadWord(addr uint16) uint16 {
  return uint16(m.ram[addr]) | uint16(m.ram[addr+1]) << 8
}

func (m *MMU) Write(addr uint16, val byte) {
  m.ram[addr] = val
}

func (m *MMU) WriteWord(addr uint16, val uint16) {
  m.ram[addr] = byte(val >> 4)
  m.ram[addr+1] = byte(val & 0b0000_1111)
}
