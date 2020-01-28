package main

import _ "fmt"

type GPU struct {
	mode      int
	modeclock int
	line      int
}

func NewGPU() *GPU {
	return &GPU{
		mode:      0,
		modeclock: 0,
		line:      0,
	}
}

func (gpu *GPU) Step() {
	switch gpu.mode {
	// OAM read mode, scanline active
	case 2:

	// VRAM read mode, scanline active
	// tread end of mode 3 as end of scanline
	case 3:

		// Hblank
		// after last hblank, push the screen data
	}
}

func (gpu *GPU) Run(delta int64) {
	// fmt.Printf("time %d\n", delta)
}
