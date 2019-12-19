package main

func main() {
  var cpu = NewCPU()
  cpu.mmu.Load("./bios.gb")

  for n := 0; n < 50; n++ {
    cpu.Step(cpu.pc)
  }
}
