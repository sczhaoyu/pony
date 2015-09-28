package server

import (
	"github.com/sczhaoyu/pony/common"
)

type LogicServerManager struct {
	MaxConn      int                   //最大链接
	ConnChan     chan *LogicConn       //链接通道
	SendChan     chan []byte           //发送数据通道
	RspChan      chan *common.Response //回应数据通道
	ClientServer *Server               //客户端链接服务器
}

func NewLogicServerManager(c *Server) *LogicServerManager {
	var ls LogicServerManager
	ls.MaxConn = 2
	ls.ConnChan = make(chan *LogicConn, ls.MaxConn)
	ls.SendChan = make(chan []byte, ls.MaxConn)
	ls.RspChan = make(chan *common.Response, ls.MaxConn)
	ls.ClientServer = c
	return &ls
}
func (l *LogicServerManager) Start() {
	for j := 0; j < l.MaxConn; j++ {
		conn := NewLogicConn(l, l.ClientServer)
		conn.Start()
		l.ConnChan <- conn
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

//链接全部重链接
func (l *LogicServerManager) ResetConnAll() {
	var ret []*LogicConn = make([]*LogicConn, 0, len(l.ConnChan))
	for v := range l.ConnChan {
		ret = append(ret, v)
		if len(ret) == len(l.ConnChan) {
			break
		}
	}
	for i := 0; i < len(ret); i++ {
		if ret[i].Conn != nil {
			ret[i].Conn.Close()
		}
		l.ReturnConn(ret[i])
	}
}
