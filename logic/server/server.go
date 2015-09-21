package server

import (
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
)

var (
	s Server
)

type Server struct {
	Ip            string              //服务器IP
	Port          int                 //启动端口
	MaxClient     int                 //服务器最大链接
	MCC           chan int            //链接处理通道
	MaxPush       int                 //推送消息最大处理数量
	RspC          chan *Response      //推送消息数据通道
	HeartbeatTime int64               //心跳超时回收时间(秒)
	MaxDataLen    int                 //最大接受数据长度
	Session       map[string]net.Conn //session
	SessionMutex  sync.Mutex          //会话操作锁
}

//创建服务
func NewServer(port int) *Server {
	s.Port = port
	s.Ip = ""
	s.MaxClient = 200
	s.MaxPush = 50000
	s.HeartbeatTime = 20
	s.MaxDataLen = 2048
	s.RspC = make(chan *Response, s.MaxPush)
	s.MCC = make(chan int, s.MaxClient)
	s.Session = make(map[string]net.Conn)
	return &s
}

//消息接受状态确认检测
func (s *Server) RspMsgCheck() {

}

//启动服务
func (s *Server) Start() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(s.Ip), s.Port, ""})
	if err != nil {
		log.Println("logic server start error:", err.Error())
		return
	}
	log.Println("logic server start success:", s.Port)
	//启动消息发送线程
	go s.sendMsg()
	//消息检测线程
	go s.RspMsgCheck()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client server error:", err.Error())
			s.CloseConn(conn)
			continue
		}
		s.MCC <- 1
		s.AddSession(conn)
		// //读取数据
		go s.ReadData(conn)
	}

}

//读取客户端服务器过来的数据
func (s *Server) ReadData(conn *net.TCPConn) {
	for {
		data, err := util.ReadData(conn, s.MaxDataLen)
		if err != nil {
			s.CloseConn(conn)
			return
		}
		//进入处理
		go handler(NewConn(conn, s), data)
	}
}

func (s *Server) AddSession(conn net.Conn) {
	s.SessionMutex.Lock()
	s.Session[conn.RemoteAddr().String()] = conn
	s.SessionMutex.Unlock()
}

//关闭链接,删除session
func (s *Server) CloseConn(conn net.Conn) {
	s.SessionMutex.Lock()
	delete(s.Session, conn.RemoteAddr().String())
	s.SessionMutex.Unlock()
	conn.Close()
	<-s.MCC
}

//加入消息
func (s *Server) Put(r *Response) {
	s.RspC <- r
}

//获取链接
func (s *Server) GetConn(r *Response) net.Conn {
	return s.Session[r.Head.Addr]
}

//发送消息
func (s *Server) sendMsg() {
	for {
		rsp := <-s.RspC
		log.Println(rsp)
		s.GetConn(rsp).Write(rsp.GetJson())
	}
}
