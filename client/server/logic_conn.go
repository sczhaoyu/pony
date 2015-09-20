package server

import (
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
	"time"
)

type LogicConn struct {
	net.Conn              //会话
	State     bool        //链接状态
	RC        chan int    //重置通道信号
	ConnMutex sync.Mutex  //数据发送锁
	DataCh    chan []byte //数据发送通道
	Addr      string      //服务器链接地址 IP+Port格式
}

//创建连接
func NewLogicConn(addr string) (*LogicConn, error) {
	var lc LogicConn
	lc.Addr = addr
	conn, err := lc.newConn(addr)
	if err != nil {

		return nil, err
	}
	lc.Conn = conn
	lc.State = true
	lc.RC = make(chan int, 1)
	lc.DataCh = make(chan []byte, 100)
	return &lc, err
}
func (lc *LogicConn) Start() {
	//状态监测
	go lc.CheckClient()
	go lc.ReadData()
}

//读取数据
func (lc *LogicConn) ReadData() {
	for {
		_, err := util.ReadData(lc.Conn)
		if err != nil {
			lc.State = false
			lc.ConnMutex.Lock()
			lc.RC <- 0
		}

	}
}

//阻塞发送数据
func (lc *LogicConn) SendData() {
	for {
		data := <-lc.DataCh
		lc.ConnMutex.Lock()
		lc.Conn.Write(data)
		lc.ConnMutex.Unlock()
		log.Println(string(data))
	}
}

//检查链接的完整
func (lc *LogicConn) CheckClient() {
	for {
		<-lc.RC
		for lc.State == false {
			if lc.Conn != nil {
				lc.Conn.Close()
			}
			var err error
			lc.Conn, err = lc.newConn(lc.Addr)
			if err == nil {
				lc.State = true
				lc.ConnMutex.Unlock()
				log.Println("logic server reset client success")
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}
}

//创建一个链接
func (lc *LogicConn) newConn(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, err
}
