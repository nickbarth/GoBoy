package main

import "time"


func main() {
  const ClockSpeed = 4194304
  const FrameRate = 60

  cpu := NewCPU()
  gpu := NewGPU()

  cpu.mmu.Load("./bios.bin")
  now := time.Now()
  cycles := 0

  for {
    delta := time.Since(now)
    now = time.Now()

    // cycles += 
    cpu.Run(delta.Microseconds())

    if cycles >= (ClockSpeed / FrameRate) {
      gpu.Run(delta.Microseconds())
    }
  }
}
