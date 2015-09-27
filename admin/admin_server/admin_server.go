package admin_server

import (
	"github.com/sczhaoyu/pony/common"
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

var admin AdminServer

//创建服务
func NewAdminServer(port int) *AdminServer {
	admin.Port = port
	admin.Ip = "127.0.0.1"
	admin.CS = make(map[string]*ClientSer)
	return &admin
}

//启动管理服务器
func (a *AdminServer) Run() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(a.Ip), a.Port, ""})
	if err != nil {
		log.Println("admin server start error:", err.Error())
		return
	}
	log.Println("admin server start success:", a.Port)
	go HtppRun()
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

func (a *AdminServer) AddSession(c net.Conn, serType, addr string, num int) {
	a.mutex.Lock()
	var cs ClientSer
	cs.Addr = addr
	cs.ServerType = serType
	cs.Conn = c
	cs.ClientNum = num
	a.CS[c.RemoteAddr().String()] = &cs
	defer a.mutex.Unlock()
}
func (a *AdminServer) Close(conn net.Conn) {
	a.mutex.Lock()
	delete(a.CS, conn.RemoteAddr().String())
	conn.Close()
	defer a.mutex.Unlock()
}
func (a *AdminServer) GetLS() []string {
	var ret []string = make([]string, 0, len(a.CS))
	for _, v := range a.CS {
		if v.ServerType == common.LS {
			ret = append(ret, v.Addr)
		}
	}
	return ret
}
