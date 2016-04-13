package server

import (
	"errors"
	"github.com/sczhaoyu/pony/common"
	"github.com/sczhaoyu/pony/util"
	// "log"
	"net"
	"sync"
	"time"
)

type LogicConn struct {
	net.Conn                         //会话
	State        bool                //链接状态
	RC           chan int            //重置通道信号
	ConnMutex    sync.Mutex          //数据发送锁
	DataCh       chan []byte         //数据发送通道
	Addr         string              //服务器链接地址 IP+Port格式
	MaxDataLen   int                 //最大接受数据长度
	LSM          *LogicServerManager //逻辑服务管理者
	ResetTimeOut int                 //超时重链接秒
	ClientServer *Server             //所属服务器
}

//创建连接
func NewLogicConn(lsm *LogicServerManager, cls *Server) *LogicConn {
	var lc LogicConn
	lc.ResetTimeOut = 2
	lc.MaxDataLen = 2048
	lc.RC = make(chan int, 1)
	lc.DataCh = make(chan []byte, 100)
	lc.LSM = lsm
	lc.ClientServer = cls
	return &lc
}
func (lc *LogicConn) Start() {
	go func() {
		conn, err := lc.newConn()
		if err != nil {
			lc.State = false
			lc.RC <- 0
		} else {
			// log.Println("logic server success:", lc.Addr)
			lc.Conn = conn
			go lc.ReadData()
		}

	}()
	//状态监测
	go lc.CheckClient()
	go lc.SendData()

}

//第一次发送数据，告诉服务器自己属于那一台客户端服务器
func (lc *LogicConn) FirstSend() {
	req := common.AuthRequest(common.LOGICCLIENT, []byte(lc.ClientServer.Listen.Addr().String()))
	lc.DataCh <- req.GetJson()
}

//读取数据
func (lc *LogicConn) ReadData() {
	lc.FirstSend()
	for {
		data, err := util.ReadData(lc.Conn, lc.MaxDataLen)
		if err != nil {
			lc.State = false
			lc.RC <- 0
			break
		}
		go handler(lc, data)
	}
}

//阻塞发送数据
func (lc *LogicConn) SendData() {
	for {
		data := <-lc.DataCh
		lc.ConnMutex.Lock()
		lc.Conn.Write(data)
		lc.ConnMutex.Unlock()
	}
}

//检查链接的完整
func (lc *LogicConn) CheckClient() {
	for {
		<-lc.RC
		lc.ConnMutex.Lock()
		for lc.State == false {
			if lc.Conn != nil {
				lc.Conn.Close()
			}
			var err error
			lc.Conn, err = lc.newConn()
			if err == nil {
				lc.State = true
				// log.Println("logic server reset client success:", lc.Conn.RemoteAddr().String())
				go lc.ReadData()
			} else {
				time.Sleep(time.Second * time.Duration(lc.ResetTimeOut))
			}
		}
		lc.ConnMutex.Unlock()
	}
}

//创建一个链接
func (lc *LogicConn) newConn() (net.Conn, error) {
	data, derr := util.HttpRequest("http://127.0.0.1:3869/logic", "post", nil, nil)
	if derr != nil {
		return nil, derr
	}
	ret := common.GetLSAddr(data)
	if ret == nil {
		return nil, errors.New("not fount server addr!")
	}
	conn, err := net.Dial("tcp", ret.Addr)
	if err != nil {
		return nil, err
	}
	lc.Addr = ret.Addr
	return conn, err
}
