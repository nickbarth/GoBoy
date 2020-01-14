package main

import _"fmt"

type GPU struct {}

func NewGPU() *GPU {
  return &GPU{}
}

func (gpu *GPU) Step() {

}

func (gpu *GPU) Run(delta int64) {
  // fmt.Printf("time %d\n", delta)
}
