package server

import (
	"github.com/sczhaoyu/pony/common"
	"github.com/sczhaoyu/pony/util"
	"time"
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
	ls.MaxConn = 50
	ls.ConnChan = make(chan *LogicConn, ls.MaxConn)
	ls.SendChan = make(chan []byte, ls.MaxConn)
	ls.RspChan = make(chan *common.Response, ls.MaxConn)
	ls.ClientServer = c
	return &ls
}
func (l *LogicServerManager) Start() {
	data, err := util.HttpRequest("http://127.0.0.1:3869/logic/list", "post", nil, nil)
	if err == nil {
		ret := common.UnmarshalLSAddr(data)
		if len(ret) == 0 {
			time.AfterFunc(time.Second*2, func() {
				l.Start()
			})
		} else {
			for i := 0; i < len(ret); i++ {
				//初始化链接
				for j := 0; j < l.MaxConn; j++ {
					conn := NewLogicConn(l)
					conn.Start()
					l.ConnChan <- conn
				}
				go l.SendLogic()
			}
		}

	} else {
		time.AfterFunc(time.Second*2, func() {
			l.Start()
		})
	}

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
