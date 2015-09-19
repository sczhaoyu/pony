package main

import (
	. "github.com/sczhaoyu/pony/client/server"
	"runtime"
)

func main() {
	NewServer(8555).Start()
	runtime.SetCPUProfileRate(runtime.NumCPU())
}
