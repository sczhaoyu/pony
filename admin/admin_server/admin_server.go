package admin_server

import (
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
)

type AdminServer struct {
	Port  int                   //启动端口
	Ip    string                //IP地址
	CS    map[string]*ClientSer //链接服务器
	mutex sync.Mutex            //session操作锁
}

//创建服务
func NewAdminServer(port int) *AdminServer {
	var a AdminServer
	a.Port = port
	a.Ip = "127.0.0.1"
	a.CS = make(map[string]*ClientSer)
	return &a
}

//启动管理服务器
func (a *AdminServer) Run() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(a.Ip), a.Port, ""})
	if err != nil {
		log.Println("admin server start error:", err.Error())
		return
	}
	log.Println("admin server start success:", a.Port)
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			a.Close(conn)
			continue
		}
		//读取数据
		go a.ReadData(conn)
	}
}

//读取客户端服务器过来的数据
func (a *AdminServer) ReadData(conn net.Conn) {
	for {
		data, err := util.ReadData(conn, 2048)
		if err != nil {
			a.Close(conn)
			return
		}
		//进入处理
		go handler(a, conn, data)
	}
}

//发送通知
func (a *AdminServer) SendNotice(st string, data []byte) {
	for _, v := range a.CS {
		if v.ServerType == st {
			v.Write(data)
		}
	}
}

func (a *AdminServer) AddSession(c net.Conn, serType string) {
	a.mutex.Lock()
	var cs ClientSer
	cs.Addr = c.RemoteAddr().String()
	cs.ServerType = serType
	cs.Conn = c
	a.CS[cs.Addr] = &cs
	defer a.mutex.Unlock()
}
func (a *AdminServer) Close(conn net.Conn) {
	a.mutex.Lock()
	delete(a.CS, conn.RemoteAddr().String())
	conn.Close()
	defer a.mutex.Unlock()
}
