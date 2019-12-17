package main

func main() {
  var cpu = NewCPU()
  cpu.mmu.Load("./bios.bin")

  for n := 0; n < 30; n++ {
    cpu.Step(cpu.pc)
  }
}
