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
	LogicSendChan  chan interface{}        //逻辑服务发送数据通道

}

//创建服务
func NewServer(port int) *Server {
	var s Server
	s.Port = port
	s.Ip = ""
	s.MaxClient = 200
	s.MaxSendLogic = 5000
	s.SessionTimeOut = 20
	s.LogicSendChan = make(chan interface{}, s.MaxSendLogic)
	s.MaxClientChan = make(chan int, s.MaxClient)
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
	go NewLogicServerConn("127.0.0.1:8456").Start()
	for {
		conn, err := listen.AcceptTCP()
		//加入会话
		s.AddSession(conn)
		if err != nil {
			log.Println("client error:", err.Error())
			s.CloseConn(conn)
			continue
		}
		s.MaxClientChan <- 1
		go s.ReadData(conn)
	}

}

//读取客户端的数据
func (s *Server) ReadData(conn *net.TCPConn) {
	for {
		var l int = 4
		data := make([]byte, l)
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(time.Second * time.Duration(s.SessionTimeOut))
			timeout <- true
		}()
		go func() {
			for l > 0 {
				i, err := conn.Read(data)
				if err != nil {
					s.CloseConn(conn)
					return
				}
				l = l - i
			}
			l = util.ByteSliceToInt(data)
			for l > 0 {
				data = make([]byte, l)
				i, err := conn.Read(data)
				if err != nil {
					s.CloseConn(conn)
					return
				}
				l = l - i
				//发送给逻辑服务器
				s.SendLogic(data)
			}
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
func (s *Server) SendLogic(r interface{}) {
	//写入数据库 备份
	s.LogicSendChan <- r
	log.Println(string(r.([]byte)))
}

//添加session
func (s *Server) AddSession(conn *net.TCPConn) {
	k := conn.RemoteAddr().String()
	s.SessionMutex.Lock()
	s.Session[k] = conn
	defer s.SessionMutex.Unlock()
}
