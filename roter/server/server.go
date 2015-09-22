package server

import (
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
)

type RoterServer struct {
	Ip           string                  //服务器IP
	Port         int                     //启动端口
	Session      map[string]*net.TCPConn //链接路由的客户端会话
	SessionMutex sync.Mutex              //会话操作锁
	LSM          *LogicServerManager     //逻辑服务管理
	MaxDataLen   int                     //最大接受数据长度
}

//创建服务
func NewRoterServer(port int) *RoterServer {
	var s RoterServer
	s.Port = port
	s.Ip = ""
	s.MaxDataLen = 2048
	s.Session = make(map[string]*net.TCPConn)
	return &s
}

//启动服务
func (s *RoterServer) Run() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(s.Ip), s.Port, ""})
	if err != nil {
		log.Println("client server start error:", err.Error())
		return
	}
	log.Println("client server start success:", s.Port)
	//启动逻辑服务器链接
	s.LSM = NewLogicServerManager("127.0.0.1:8456")
	go s.LSM.Start()
	go s.ResponseHandle()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client error:", err.Error())
			s.CloseConn(conn)
			continue
		}
		go s.ReadData(conn)
	}

}

//读取客户端的数据
func (s *RoterServer) ReadData(conn net.Conn) {
	for {
		data, err := util.ReadData(conn, s.MaxDataLen)
		if err != nil {
			s.CloseConn(conn)
			return
		}
		//提交给逻辑服务器
		go s.SendLogic(data)

	}
}

//发送逻辑服务器处理
func (s *RoterServer) SendLogic(data []byte) {
	tmp := util.IntToByteSlice(len(data))
	tmp = append(tmp, data...)
	//写入数据库 备份
	s.LSM.SendChan <- tmp
}

//关闭链接
func (s *RoterServer) CloseConn(conn net.Conn) {
	s.SessionMutex.Lock()
	//删除session
	delete(s.Session, conn.RemoteAddr().String())
	s.SessionMutex.Unlock()
	conn.Close()
}

//接收逻辑服务器回应结果
//分发处理
//缓冲区阻塞处理
func (s *RoterServer) ResponseHandle() {
	for {
		data := <-s.LSM.RspChan
		handler(s, data)
	}
}
