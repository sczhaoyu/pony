package server

import (
	"github.com/sczhaoyu/pony/common"
)

type LogicServerManager struct {
	Addr         string                //IP+Port格式
	MaxConn      int                   //最大链接
	ConnChan     chan *LogicConn       //链接通道
	SendChan     chan []byte           //发送数据通道
	RspChan      chan *common.Response //回应数据通道
	ClientServer *Server               //客户端链接服务器
}

func NewLogicServerManager(addr string, c *Server) *LogicServerManager {
	var ls LogicServerManager
	ls.Addr = addr
	ls.MaxConn = 2
	ls.ConnChan = make(chan *LogicConn, ls.MaxConn)
	ls.SendChan = make(chan []byte, ls.MaxConn)
	ls.RspChan = make(chan *common.Response, ls.MaxConn)
	ls.ClientServer = c
	return &ls
}
func (l *LogicServerManager) Start() {
	//初始化链接
	for i := 0; i < l.MaxConn; i++ {
		conn := NewLogicConn(l.Addr, l)
		l.ConnChan <- conn
		conn.Start()
	}
	go l.SendLogic()
}

//发送给已经链接的客户端
func (l *LogicServerManager) SendLogic() {
	for {
		data := <-l.SendChan

		conn := l.GetConn()
		conn.DataCh <- data
		l.ReturnConn(conn)
	}
}

//获取链接
func (l *LogicServerManager) GetConn() *LogicConn {
	ret := <-l.ConnChan
	return ret
}

//放入链接
func (l *LogicServerManager) ReturnConn(conn *LogicConn) {
	l.ConnChan <- conn
}
