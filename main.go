package main

import _ "time"

func main() {
	const ClockSpeed = 4194304
	const FrameRate = 60

	cpu := NewCPU()
	gpu := NewGPU()

	cpu.mmu.Load("./bios.bin")

	steps := 41
	steps = 20

	for i := 0; i < steps; i++ {
		gpu.Step()
		cpu.Step(cpu.pc)
	}

	cpu.mmu.Print()

	/*
		cpu.Step(cpu.pc)
		gpu.Step()

		  now := time.Now()
		  cycles := 0

		  for i := 0; i < 80; i++ {
		    delta := time.Since(now)
		    now = time.Now()
		    cpu.Run(delta.Microseconds())

		    if cycles >= (ClockSpeed / FrameRate) {
		      gpu.Run(delta.Microseconds())
		    }
		  }
	*/
}
