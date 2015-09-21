package server

import (
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip             string                  //服务器IP
	Port           int                     //启动端口
	Session        map[string]*net.TCPConn //客户端链接会话
	SessionMutex   sync.Mutex              //会话操作锁
	SessionTimeOut int64                   //会话无动作超时
	MaxClient      int                     //服务器最大链接
	MaxClientChan  chan int                //链接处理通道
	MaxSendLogic   int                     //推送客户端消息最大处理数量
	LSM            *LogicServerManager     //逻辑服务管理
	MaxDataLen     int                     //最大接受数据长度
}

//创建服务
func NewServer(port int) *Server {
	var s Server
	s.Port = port
	s.Ip = ""
	s.MaxClient = 200
	s.MaxSendLogic = 5000
	s.SessionTimeOut = 200
	s.MaxClientChan = make(chan int, s.MaxClient)
	s.MaxDataLen = 2048
	s.Session = make(map[string]*net.TCPConn)
	return &s
}

//启动服务
func (s *Server) Start() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(s.Ip), s.Port, ""})
	if err != nil {
		log.Println("client server start error:", err.Error())
		return
	}
	log.Println("client server start success:", s.Port)
	//启动逻辑服务器链接
	s.LSM = NewLogicServerManager("127.0.0.1:8456")
	go s.LSM.Start()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client error:", err.Error())
			s.CloseConn(conn)
			continue
		}
		s.MaxClientChan <- 1
		//加入会话
		s.AddSession(conn)
		go s.ReadData(conn)
	}

}

//读取客户端的数据
func (s *Server) ReadData(conn *net.TCPConn) {
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(time.Second * time.Duration(s.SessionTimeOut))
			timeout <- true
		}()
		go func() {
			data, err := util.ReadData(conn, s.MaxDataLen)
			if err != nil {
				s.CloseConn(conn)
				return
			}
			//发送给逻辑服务器
			s.SendLogic(data)
		}()
		<-timeout //超时关闭链接
		s.CloseConn(conn)
		return

	}

}

//关闭链接
func (s *Server) CloseConn(conn *net.TCPConn) {
	s.SessionMutex.Lock()
	//删除session
	delete(s.Session, conn.RemoteAddr().String())
	s.SessionMutex.Unlock()
	conn.Close()
	<-s.MaxClientChan
}

//发送逻辑服务器处理
func (s *Server) SendLogic(data []byte) {
	//写入数据库 备份
	s.LSM.SendChan <- util.ByteLen(data)
}

//添加session
func (s *Server) AddSession(conn *net.TCPConn) {
	k := conn.RemoteAddr().String()
	s.SessionMutex.Lock()
	s.Session[k] = conn
	defer s.SessionMutex.Unlock()
}
