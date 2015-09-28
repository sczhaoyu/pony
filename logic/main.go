package main

import (
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
		s.Radio([]byte("广播"))
	}()
	s.Start()
}
