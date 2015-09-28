package main

import (
	"github.com/sczhaoyu/pony/common"
	. "github.com/sczhaoyu/pony/logic/server"
	"log"
	"runtime"
	"time"
)

func main() {
	runtime.SetCPUProfileRate(runtime.NumCPU())
	intF := func(c *Conn) bool {
		log.Println("拦截器已经运行")
		return true
	}
	BeforeInterceptor(intF)
	ReadFunc = func(c *Conn) {
		log.Println("body：", string(c.Request.Body))
	}
	s := NewServer(9862)
	go func() {
		time.Sleep(time.Second * 10)
		var l common.LSAddr
		l.Addr = "12222.000,00"
		l.Num = 17666
		s.Radio(&l)
	}()
	s.Start()
}
