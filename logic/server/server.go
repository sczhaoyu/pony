package server

import (
	"github.com/sczhaoyu/pony/common"
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
)

var (
	s Server
)

type Server struct {
	Ip         string            //服务器IP
	Port       int               //启动端口
	MaxClient  int               //服务器最大链接
	MCC        chan int          //链接处理通道
	MaxPush    int               //推送消息最大处理数量
	RspC       chan *Write       //推送消息数据通道
	MaxDataLen int               //最大接受数据长度
	Session    SessionManager    //会话管理
	AdminConn  *common.AdminConn //总控制台
	Listen     *net.TCPListener
}

//创建服务
func NewServer(port int) *Server {
	s.Port = port
	s.Ip = "127.0.0.1"
	s.MaxClient = 600
	s.MaxPush = 50000
	s.MaxDataLen = 2048
	s.RspC = make(chan *Write, s.MaxPush)
	s.MCC = make(chan int, s.MaxClient)
	s.Session.Init()
	return &s
}

//启动服务
func (s *Server) Start() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(s.Ip), s.Port, ""})
	s.Listen = listen
	if err != nil {
		log.Println("logic server start error:", err.Error())
		return
	}
	log.Println("logic server start success:", s.Port)
	//启动后台管理服务器链接
	s.ConnectAdmin()
	//启动消息发送线程
	go s.sendMsg()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client server error:", err.Error())
			s.CloseConn(conn)
			continue
		}
		s.MCC <- 1
		// //读取数据
		go s.ReadData(conn)
	}

}

//读取客户端服务器过来的数据
func (s *Server) ReadData(conn *net.TCPConn) {
	//加入session
	s.Session.SetSession(conn, "")
	// //通知admin
	// s.NoticeAdmin()
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

// //通知管理服务器有新的session
// func (s *Server) NoticeAdmin() {
// 	var la common.LSAddr
// 	la.Addr = s.Listen.Addr().String()
// 	la.Num = s.Session.GetLen()
// 	req := common.AuthRequest(common.ADDLSCONN, la.Marshal())
// 	s.AdminConn.DataCh <- req.GetJson()
// }

//链接管理服务器
func (s *Server) ConnectAdmin() {
	//启动后台管理服务器链接
	s.AdminConn = common.NewAdminConn("127.0.0.1:2058")
	//初始化发送信息
	s.AdminConn.FirstSendAdmin = func() {
		req := common.AutoLSAddrReq(common.LS, s.Listen.Addr().String(), s.Session.GetLen())
		s.AdminConn.DataCh <- req.GetJson()
	}
	s.AdminConn.InitSendFunc = func() {
		req := common.AutoLSAddrReq(common.ADDLS, s.Listen.Addr().String(), s.Session.GetLen())
		s.AdminConn.DataCh <- req.GetJson()

	}
	s.AdminConn.Run()
}

//关闭链接,删除session
func (s *Server) CloseConn(conn net.Conn) {
	s.Session.RemoveAddrSession(conn.RemoteAddr().String())
	<-s.MCC
	//通知链接减去
	var la common.LSAddr
	la.Addr = s.Listen.Addr().String()
	la.Num = s.Session.GetLen()
	req := common.AuthRequest(common.DELLSCONN, la.Marshal())
	s.AdminConn.DataCh <- req.GetJson()
}

//加入消息
func (s *Server) Put(r *Write) {
	s.RspC <- r
}

//发送消息
func (s *Server) sendMsg() {
	for {
		rsp := <-s.RspC
		rsp.Out()
	}
}

//全局广播
func (s *Server) Radio(data interface{}) {
	for _, v := range s.Session.SCName {
		for _, c := range v {
			//通知前端的每台clientServer
			var rsp common.Response
			rsp.Head = new(common.ResponseHead)
			rsp.Head.Command = common.RADIO
			rsp.Body = data
			var w Write
			w.Conn = c
			w.Body = rsp.GetJson()
			s.Put(&w)
			break
		}
	}
}
