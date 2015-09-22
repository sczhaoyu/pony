package main

import (
	. "github.com/sczhaoyu/pony/roter/server"
	"runtime"
)

func main() {
	NewRoterServer(8061).Run()
	runtime.SetCPUProfileRate(runtime.NumCPU())
}
