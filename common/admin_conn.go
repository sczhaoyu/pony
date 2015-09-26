package common

import (
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
	"time"
)

type AdminConn struct {
	net.Conn                   //会话
	State          bool        //链接状态
	RC             chan int    //重置通道信号
	ConnMutex      sync.Mutex  //数据发送锁
	DataCh         chan []byte //数据发送通道
	Addr           string      //服务器链接地址 IP+Port格式
	MaxDataLen     int         //最大接受数据长度
	ResetTimeOut   int         //超时重链接秒
	FirstSendAdmin func()      //初始化发送信息
}

//创建连接
func NewAdminConn(addr string) *AdminConn {
	var a AdminConn
	a.Addr = addr
	a.ResetTimeOut = 2
	a.State = true
	a.MaxDataLen = 2048
	a.RC = make(chan int, 1)
	a.DataCh = make(chan []byte, 100)
	return &a
}
func (a *AdminConn) Run() {
	go func() {
		conn, err := a.newConn(a.Addr)
		if err != nil {
			a.State = false
			a.RC <- 0
		} else {
			log.Println("admin server success!")
			a.Conn = conn
			go a.ReadData()
		}

	}()
	//状态监测
	go a.CheckClient()
	go a.SendData()

}

//读取数据
func (a *AdminConn) ReadData() {
	go a.FirstSendAdmin()
	for {
		data, err := util.ReadData(a.Conn, a.MaxDataLen)
		if err != nil {
			a.State = false
			a.RC <- 0
			break
		}
		log.Println("admin报文:", string(data))
	}
}

//阻塞发送数据
func (a *AdminConn) SendData() {
	for {
		data := <-a.DataCh
		a.ConnMutex.Lock()
		a.Conn.Write(data)
		a.ConnMutex.Unlock()
	}
}

//检查链接的完整
func (a *AdminConn) CheckClient() {
	for {
		<-a.RC
		a.ConnMutex.Lock()
		for a.State == false {
			if a.Conn != nil {
				a.Conn.Close()
			}
			var err error
			a.Conn, err = a.newConn(a.Addr)
			if err == nil {
				a.State = true
				log.Println("admin server reset client success:", a.Conn.RemoteAddr().String())
				go a.ReadData()
			} else {
				time.Sleep(time.Second * time.Duration(a.ResetTimeOut))
			}
		}
		a.ConnMutex.Unlock()
	}
}

//创建一个链接
func (a *AdminConn) newConn(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, err
}