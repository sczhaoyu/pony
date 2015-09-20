package server

import (
	"log"
)

type LogicServer struct {
	Addr     string          //IP+Port格式
	MaxConn  int             //最大链接
	ConnChan chan *LogicConn //链接通道
}

func NewLogicServerConn(addr string) *LogicServer {
	var ls LogicServer
	ls.Addr = addr
	ls.MaxConn = 1
	ls.ConnChan = make(chan *LogicConn, ls.MaxConn)
	return &ls
}
func (l *LogicServer) Start() {
	//初始化链接
	for i := 0; i < l.MaxConn; i++ {
		conn, err := NewLogicConn(l.Addr)
		if err != nil {
			log.Println("logic server error:", err.Error())
			return
		}
		l.ConnChan <- conn
		conn.Start()
	}
	log.Println("logic server success:", l.Addr)
}

//获取链接
func (l *LogicServer) GetConn() *LogicConn {
	ret := <-l.ConnChan
	return ret
}

//放入链接
func (l *LogicServer) ReturnConn(conn *LogicConn) {
	l.ConnChan <- conn
}
