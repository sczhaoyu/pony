package server

import (
	"encoding/json"
	"errors"
	"github.com/sczhaoyu/pony/client/biz"
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip             string             //服务器IP
	Port           int                //启动端口
	Session        map[int64]net.Conn //客户端链接会话
	SessionMutex   sync.Mutex         //会话操作锁
	SessionTimeOut int64              //会话无动作超时
	MaxClient      int                //服务器最大链接
	MaxClientChan  chan int           //链接处理通道
	MaxSendLogic   int                //推送客户端消息最大处理数量
	MaxDataLen     int                //最大接受数据长度
	RSC            chan []byte        //回应客户端数据通道
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
	s.Session = make(map[int64]net.Conn)
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

		cErr, uid := s.CheckToken(conn)
		if cErr != nil {
			conn.Close()
			continue
		}
		//加入会话
		s.AddSession(conn, uid)
		s.MaxClientChan <- 1
		go s.ReadData(conn, uid)
		go s.RSCSend()
	}

}

//读取客户端的数据
func (s *Server) ReadData(conn net.Conn, uid int64) {
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(time.Second * time.Duration(s.SessionTimeOut))
			timeout <- true
		}()
		go func() {
			data, err := util.ReadData(conn, s.MaxDataLen)
			if err != nil {
				s.CloseConn(uid)
				return
			}
			//让路由器链接发送给路由器处理
			s.Roter.DataCh <- util.ByteLen(data)

		}()
		<-timeout //超时关闭链接
		s.CloseConn(uid)
		return
	}

}

//检查token
func (s *Server) CheckToken(conn net.Conn) (error, int64) {
	//直接检查是否是注册用户
	data, err := util.ReadData(conn, s.MaxDataLen)
	if err != nil {
		return err, 0
	}
	var rq Request
	err = rq.Unmarshal(data)
	if err != nil {
		return err, 0
	}
	b := biz.CheckToken(rq.Head.Token)
	if b == false {
		return errors.New("token not found"), 0
	}
	return nil, rq.Head.UserId
}

//关闭链接
func (s *Server) CloseConn(uid int64) {
	s.SessionMutex.Lock()
	s.Session[uid].Close()
	delete(s.Session, uid)
	s.SessionMutex.Unlock()
	<-s.MaxClientChan
}

//回应客户端数据
func (s *Server) RSCSend() {
	for {
		data := <-s.RSC
		var r Request
		json.Unmarshal(data, &r)
		conn := s.GetSession(r.Head.UserId)
		if conn != nil {
			conn.Write(r.GetJson())
		}
	}
}
func (s *Server) GetSession(userId int64) net.Conn {
	return s.Session[userId]
}

//添加session
func (s *Server) AddSession(conn net.Conn, uid int64) {
	s.SessionMutex.Lock()
	s.Session[uid] = conn
	defer s.SessionMutex.Unlock()
}
