package server

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/common"
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip             string                //服务器IP
	Port           int                   //启动端口
	Session        common.SessionManager //客户端链接会话
	SessionMutex   sync.Mutex            //会话操作锁
	SessionTimeOut int64                 //会话无动作超时
	MaxClient      int                   //服务器最大链接
	MaxClientChan  chan int              //链接处理通道
	MaxSendLogic   int                   //推送客户端消息最大处理数量
	MaxDataLen     int                   //最大接受数据长度
	RSC            chan []byte           //回应客户端数据通道
	Roter          *RoterConn
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
	s.RSC = make(chan []byte, s.MaxClient)
	s.Session.Init()
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
	//启动路由器链接
	r, rr := NewRoterConn("127.0.0.1:8061", s)
	if rr != nil {
		log.Println("roter server error:", rr.Error())
		return
	} else {
		log.Println("roter server success:")
	}
	s.Roter = r
	go r.Start()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client error:", err.Error())
			conn.Close()
			continue
		}
		s.MaxClientChan <- 1
		c := common.NewConn(conn)
		//加入会话
		s.AddSession(c)
		go s.ReadData(c)
		go s.RSCSend()
	}

}

//读取客户端的数据
func (s *Server) ReadData(conn *common.Conn) {
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
			data = NewRequest(conn, data).GetJson()
			//让路由器链接发送给路由器处理
			s.Roter.DataCh <- data

		}()
		<-timeout //超时关闭链接
		s.CloseConn(conn)
		return
	}

}

//关闭链接
func (s *Server) CloseConn(conn *common.Conn) {
	s.SessionMutex.Lock()
	//删除session
	delete(s.Session.Session, conn.UUID)
	s.SessionMutex.Unlock()
	conn.Close()
	<-s.MaxClientChan
}

//回应客户端数据
func (s *Server) RSCSend() {
	for {
		data := <-s.RSC
		var r Request
		json.Unmarshal(data, &r)
		conn := s.GetSession(r.Head.Cid).ClientConn
		if conn != nil {
			conn.Write(r.GetJson())
		}
	}
}
func (s *Server) GetSession(k string) *common.UsersSession {
	return s.Session.Session[k]
}

//添加session
func (s *Server) AddSession(conn *common.Conn) {
	var u common.UsersSession
	u.ClientConn = conn
	u.UserAddr = conn.RemoteAddr().String()
	u.UUID = conn.UUID
	s.Session.Session[u.UUID] = &u
	defer s.SessionMutex.Unlock()
}
