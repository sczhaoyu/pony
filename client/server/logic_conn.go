package server

import (
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
	"time"
)

type LogicConn struct {
	net.Conn                       //会话
	State      bool                //链接状态
	RC         chan int            //重置通道信号
	ConnMutex  sync.Mutex          //数据发送锁
	DataCh     chan []byte         //数据发送通道
	Addr       string              //服务器链接地址 IP+Port格式
	MaxDataLen int                 //最大接受数据长度
	LSM        *LogicServerManager //逻辑服务管理者
}

//创建连接
func NewLogicConn(addr string, lsm *LogicServerManager) (*LogicConn, error) {
	var lc LogicConn
	lc.Addr = addr
	conn, err := lc.newConn(addr)
	if err != nil {

		return nil, err
	}
	lc.Conn = conn
	lc.State = true
	lc.MaxDataLen = 2048
	lc.RC = make(chan int, 1)
	lc.DataCh = make(chan []byte, 100)
	lc.LSM = lsm
	return &lc, err
}
func (lc *LogicConn) Start() {
	//状态监测
	go lc.CheckClient()
	go lc.ReadData()
	go lc.SendData()

}

//读取数据
func (lc *LogicConn) ReadData() {
	for {
		data, err := util.ReadData(lc.Conn, lc.MaxDataLen)
		if err != nil {
			lc.State = false
			lc.RC <- 0
			break
		}
		log.Println(string(data))
		lc.LSM.RspChan <- data
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
			lc.Conn, err = lc.newConn(lc.Addr)
			if err == nil {
				lc.State = true
				log.Println("roter server reset client success")
				go lc.ReadData()
			} else {
				time.Sleep(time.Second * 5)
			}
		}
		lc.ConnMutex.Unlock()
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
