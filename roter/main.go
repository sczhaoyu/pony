package main

import (
	. "github.com/sczhaoyu/pony/roter/server"
	"runtime"
)

func main() {
	runtime.SetCPUProfileRate(runtime.NumCPU())
	NewRoterServer(8061).Run()

}
