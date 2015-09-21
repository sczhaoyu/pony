package main

import (
	. "github.com/sczhaoyu/pony/logic/server"
	"runtime"
)

func main() {
	runtime.SetCPUProfileRate(runtime.NumCPU())
	NewServer(8456).Start()
}
