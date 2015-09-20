package server

import (
	"log"
	"net"
	"sync"
	"time"
)

var (
	s Server
)

type Server struct {
	Ip             string             //服务器IP
	Port           int                //启动端口
	MaxClient      int                //服务器最大链接
	MaxClientChan  chan int           //链接处理通道
	MaxResponds    int                //推送消息最大处理数量
	RespondsChan   chan *Response     //推送消息数据通道
	HeartbeatTime  int64              //心跳超时回收时间(秒)
	RspMsg         map[string]*RspMsg //已发送的消息
	RspSendTimeOut int64              //等待发送回应超时(秒)
	APMutex        sync.Mutex         //rspMsg and pushSendQueue Mutex
}

//创建服务
func NewServer(port int) *Server {

	s.Port = port
	s.Ip = ""
	s.MaxClient = 200
	s.MaxResponds = 50000
	s.HeartbeatTime = 20
	s.RspSendTimeOut = 180
	s.RespondsChan = make(chan *Response, s.MaxResponds)
	s.MaxClientChan = make(chan int, s.MaxClient)
	s.RspMsg = make(map[string]*RspMsg)
	return &s
}

//消息接受状态确认检测
func (s *Server) RspMsgCheck() {
	for _, v := range s.RspMsg {
		if v.State == false && time.Now().Unix()-s.RspSendTimeOut > v.SendTime {
			s.pushSendQueue(v.Req, v.Rsp, true)
		}
	}
	time.Sleep(time.Second * time.Duration(s.RspSendTimeOut))
}

//启动服务
func (s *Server) Start() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(s.Ip), s.Port, ""})
	if err != nil {
		log.Println("logic server start error:", err.Error())
		return
	}
	log.Println("logic server start success:", s.Port)
	// //读取数据
	// go s.ReadData(conn)
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
		s.MaxClientChan <- 1
	}

}

//读取客户端服务器过来的数据
func (s *Server) ReadData(conn *net.TCPConn) {
	data := make([]byte, 4)
	for {
		_, err := conn.Read(data)
		if err != nil {
			s.CloseConn(conn)
			break
		}
		//进入路由器
		go handler(data)
	}

}

//关闭链接
func (s *Server) CloseConn(conn *net.TCPConn) {
	conn.Close()
	<-s.MaxClientChan
}

//加入消息
func (s *Server) AppendRspMsg(r *Request, data interface{}) {

	var msg RspMsg
	msg.Req = r
	msg.Rsp = data
	msg.SendTime = time.Now().Unix()
	msg.State = false
	s.RspMsg[r.Head.Uuid] = &msg

}

//获取链接
func (s *Server) GetConn(r *Response) *net.TCPConn {
	return nil
}

//将对象写入到发送的队列
func (s *Server) pushSendQueue(r *Request, data interface{}, b bool) {
	s.APMutex.Lock()
	rsp := NewClientResponse(r, data)
	if b {
		//更新数据库消息ID的时间
		//更新重发队列的时间
		s.RspMsg[r.Head.Uuid].SendTime = time.Now().Unix()
	} else {
		//写入等待重发确认
		s.AppendRspMsg(r, data)
		//写入数据库备份
	}
	s.RespondsChan <- rsp
	defer s.APMutex.Unlock()
}

//发送消息
func (s *Server) sendMsg() {
	for {
		rsp := <-s.RespondsChan
		s.GetConn(rsp).Write(rsp.GetJson())
	}
}
