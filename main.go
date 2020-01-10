package main

import "time"

func main() {
  // cpu := NewCPU()
  gpu := NewGPU()

  // cpu.mmu.Load("./bios.bin")

  now := time.Now()

  for {
    delta := time.Since(now)
    now = time.Now()
    gpu.Run(delta.Microseconds())
  }
}
