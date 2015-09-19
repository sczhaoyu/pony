package server

import (
	"log"
	"net"
)

type Server struct {
	Ip            string            //服务器IP
	Port          int               //启动端口
	MaxClient     int               //服务器最大链接
	MaxClientChan chan *net.TCPConn //链接处理通道
	MaxSendLogic  int               //推送客户端消息最大处理数量
	LogicSendChan chan interface{}  //逻辑服务发送数据通道
	HeartbeatTime int64             //心跳超时回收时间(秒)

}

//创建服务
func NewServer(port int) *Server {
	var s Server
	s.Port = port
	s.Ip = ""
	s.MaxClient = 1
	s.MaxSendLogic = 5000
	s.HeartbeatTime = 8000
	s.LogicSendChan = make(chan interface{}, s.MaxSendLogic)
	s.MaxClientChan = make(chan *net.TCPConn, s.MaxClient)
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
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client error:", err.Error())
			s.CloseConn(conn)
			continue
		}
		s.MaxClientChan <- conn
		go s.ReadData(conn)
	}

}
func (s *Server) ReadData(conn *net.TCPConn) {
	data := make([]byte, 158)
	for {
		_, err := conn.Read(data)
		if err != nil {
			s.CloseConn(conn)
			break
		}
		//发送给逻辑服务器
		s.SendLogic(data)
	}

}

//关闭链接
func (s *Server) CloseConn(conn *net.TCPConn) {
	conn.Close()
	<-s.MaxClientChan
}

//发送逻辑服务器处理
func (s *Server) SendLogic(r interface{}) {
	//写入数据库 备份
	s.LogicSendChan <- r
	log.Println(string(r.([]byte)))
}
