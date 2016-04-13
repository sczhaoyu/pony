package main

import (
	. "github.com/sczhaoyu/pony/business/srv"
	"runtime"
)

//该模块主要实现业务逻辑
//可以启动多个
func main() {
	var server BusinessServer
	runtime.GOMAXPROCS(runtime.NumCPU())
	server.RunCustomerServer()
	server.RunBusinessServer()

}
