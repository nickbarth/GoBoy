package main

func main() {
  var cpu = NewCPU()
  cpu.mmu.Load("./test.gb")

  for n := 0; n < 10; n++ {
    cpu.Step(n)
  }
}
