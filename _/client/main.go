package main

import (
	. "github.com/sczhaoyu/pony/client/server"
	"runtime"
)

func main() {
	runtime.SetCPUProfileRate(runtime.NumCPU())
	NewServer(1587).Start()

}
