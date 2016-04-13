package main

import (
	. "github.com/sczhaoyu/pony/admin/srv"
	"runtime"
)

//该模块主要管理全部模块
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Run()
}
