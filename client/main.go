package main

import (
	. "github.com/sczhaoyu/pony/client/srv"
	"runtime"
)

//该模块主要管理管理客户端链接
func main() {
	var c ClientServer
	runtime.GOMAXPROCS(runtime.NumCPU())
	c.RunCustomerServer()
	c.RunSrv()

}
