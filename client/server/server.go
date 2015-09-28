package server

import (
	"github.com/sczhaoyu/pony/common"
	"github.com/sczhaoyu/pony/session"
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip             string          //服务器IP
	Port           int             //启动端口
	Session        session.Manager //客户端链接会话
	SessionMutex   sync.Mutex      //会话操作锁
	SessionTimeOut int64           //会话无动作超时
	MaxClient      int             //服务器最大链接
	MaxClientChan  chan int        //链接处理通道
	MaxSendLogic   int             //推送客户端消息最大处理数量
	MaxDataLen     int             //最大接受数据长度
	RSC            chan []byte     //回应客户端数据通道
	LSM            *LogicServerManager
	Listen         *net.TCPListener
}

//创建服务
func NewServer(port int) *Server {
	var s Server
	s.Port = port
	s.Ip = "127.0.0.1"
	s.MaxClient = 200
	s.MaxSendLogic = 5000
	s.SessionTimeOut = 200
	s.MaxClientChan = make(chan int, s.MaxClient)
	s.MaxDataLen = 2048
	s.Session.Init()
	return &s
}

//启动服务
func (s *Server) Start() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(s.Ip), s.Port, ""})
	s.Listen = listen
	if err != nil {
		log.Println("client server start error:", err.Error())
		return
	}
	log.Println("client server start success:", listen.Addr().String())
	go s.AdminConnRun()
	//启动逻辑服务器链接
	lsm := NewLogicServerManager(s)
	s.LSM = lsm
	go s.LSM.Start()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client error:", err.Error())
			conn.Close()
			continue
		}
		//加入会话
		s.MaxClientChan <- 1
		go s.ReadData(s.Session.SetSession(conn))
		go s.RSCSend()
	}

}
func (s *Server) AdminConnRun() {
	//启动后台管理服务器链接
	a := common.NewAdminConn("127.0.0.1:2058")
	a.FirstSendAdmin = func() {
		//获取链接
		rsp := common.AuthResponse(common.GETLS, []byte(" "))
		//登记自己
		lg := common.AuthResponse(common.CS, s.Listen.Addr().String())
		a.DataCh <- lg.GetJson()
		a.DataCh <- rsp.GetJson()
	}
	a.RspFunc = s.AdminRsp
	a.Run()
}

//收到通知让全部链接重连
func (s *Server) AdminRsp(rsp *common.Response) {
	if rsp.Head.Command == common.ADDLS {
		log.Println("loading reset logic server client .......")
		s.LSM.ResetConnAll()
		log.Println("loading reset logic server client success!")
	}

}

//读取客户端的数据
func (s *Server) ReadData(session *session.Session) {
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(time.Second * time.Duration(s.SessionTimeOut))
			timeout <- true
		}()
		go func() {
			data, err := util.ReadData(session, s.MaxDataLen)
			if err != nil {
				s.CloseConn(session)
				return
			}

			//转发逻辑服务端链接
			s.LSM.SendChan <- common.NewRequestJson(data, session.SESSIONID)

		}()
		<-timeout //超时关闭链接
		s.CloseConn(session)
		return
	}

}

//关闭链接
func (s *Server) CloseConn(si *session.Session) {
	s.Session.Remove(si)
	<-s.MaxClientChan
}

//回应客户端数据
func (s *Server) RSCSend() {
	for {
		rsp := <-s.LSM.RspChan
		conn := s.Session.GetSession(rsp.SessionId)
		if conn != nil {
			conn.Write(rsp.Data)
		}

	}
}
