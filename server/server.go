package server

import (
	"fmt"
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"strconv"
	"time"
)

//该类主要用于服务器端的链接
//接收和回应数据包，前4字节为数据包大小
type Server struct {
	Name             string          //服务器模块名称
	Id               string          //服务器的唯一ID信息
	Ip               string          //服务器IP
	Port             int             //启动端口
	TimeOut          int64           //会话无动作超时
	MaxClient        int             //服务器最大链接
	maxClientChan    chan int        //链接处理通道
	MaxData          int             //最大接受数据长度（字节）
	Handle           func(*Conn)     //处理信息接口
	*net.TCPListener                 //TCP链接信息
	SessionTimeOut   int64           //session会话超时时间(秒)
	DPM              *DataPkgManager //数据包管理
	MemProvider      *MemProvider    //session管理工具
}

//数据包重发
func (s *Server) rewire() {
	for {
		pkg := <-s.DPM.CH
		if ok := s.MemProvider.SessionExist(pkg.Addr); ok {
			//发送数据包，获取session
			session, _ := s.MemProvider.SessionRead(pkg.Addr)
			var rsp Respon
			rsp.Unmarshal(pkg.Data[0:4])
			//重发数据包
			i, err := session.Write(pkg.Data)
			if err != nil && i != len(pkg.Data) {
				//说明数据写入失败，入库保存
			} else {
				//重新加入，等待重发确认
				s.DPM.AddPkg(rsp.Header.ResponsId, session.RemoteAddr().String(), rsp.Header.UserUid, pkg.Data)
			}
		}
	}
}
func (s *Server) init() {
	if s.SessionTimeOut == 0 {
		s.SessionTimeOut = 20
	}
	s.DPM = NewDataPkgManager(5, 5)
	if s.Id == "" {
		s.Id = util.GetUUID()
	}
	if s.TimeOut == 0 {
		//默认300秒,没有任何数据就关闭链接
		s.TimeOut = 300
	}
	if s.MaxClient == 0 {
		//默认最大100个链接
		s.MaxClient = 10
	}
	if s.MaxData == 0 {
		//未设置读取报文大小，默认800字节
		s.MaxData = 800
	}
	//判断启动的端口是否为空，如果为空给一个默认端口
	if s.Port == 0 {
		//默认8888端口启动
		s.Port = 8888
	}
	//指定session会话的超时时间
	s.MemProvider = NewSessionStore(s.SessionTimeOut)
	//指定最大链接数量
	s.maxClientChan = make(chan int, s.MaxClient)
	if s.Name == "" {
		//服务器名字没有指定使用IP加端口为服务器名称
		s.Name = s.Ip + ":" + strconv.Itoa(s.Port)
	}
}
func (s *Server) GetSessionID(conn net.Conn) string {
	sid := conn.RemoteAddr().String()
	return sid
}

//启动服务
func (s *Server) Run() {
	//调用初始数据
	s.init()
	//创建链接
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(s.Ip), s.Port, ""})
	//设置服务的链接信息
	s.TCPListener = listen
	if err != nil {
		log.Println(fmt.Sprintf("[%s]start error:%d", s.Name, s.Port))
		return
	}
	log.Println(fmt.Sprintf("[%s]start success port:%d", s.Name, s.Port))
	//线程重发
	go s.rewire()
	for {
		conn, err := s.Accept()
		//加入session
		s.addSession(conn)
		if err != nil {
			s.Close(conn)
			continue
		}
		//超时监控
		timeout := make(chan bool, 1)
		//超过超时时间，往超时通道写入一条信息
		if s.TimeOut > 0 {
			go func() {
				time.Sleep(time.Second * time.Duration(10))
				timeout <- true
			}()
		}
		select {
		//加入会话
		case s.maxClientChan <- 1:
			go s.Read(conn)
		case <-timeout:
			s.Close(conn)
		}
	}
}
func (s *Server) addSession(conn net.Conn) {
	//唯一的sessionID编号
	sid := s.GetSessionID(conn)
	session, _ := s.MemProvider.SessionRead(sid)
	session.Conn, _ = NewConn(conn, s, session, nil)

}

//读取客户端的数据
func (s *Server) Read(conn net.Conn) {
	for {
		//超时读取通道
		timeOut := make(chan bool, 1)
		timeRead := make(chan bool, 1)
		//超过超时时间，往超时通道写入一条信息
		go func() {
			time.Sleep(time.Second * time.Duration(s.TimeOut))
			timeOut <- true
		}()
		//读取信息,如果在规定时间内还没有读取完，将关闭链接
		go func() {
			data, err := util.ReadData(conn, s.MaxData)
			if err != nil {
				timeOut <- true
				return
			}
			//正常读取
			timeRead <- true
			if s.Handle != nil {
				if c := s.handleData(conn, data); c != nil {
					s.Handle(c)
				}

			}
		}()
		select {
		case <-timeRead: //正常读取
			continue
		case <-timeOut: //超时读取
			s.Close(conn)
			return
		}

	}

}

//处理数据包
func (s *Server) handleData(conn net.Conn, data []byte) *Conn {
	session, _ := s.MemProvider.SessionRead(s.GetSessionID(conn))
	c, err := NewConn(conn, s, session, data)
	//说明数据解析不正确
	if err != nil {
		c.WriteJson(err)
		//不做处理，回应客户端错误消息
		return nil
	}
	//如果该数据包处理过，直接跳过。
	ok := s.DPM.Receive(c.Request.Header.RequestId)
	if ok {
		log.Println("continue data pkg id:", c.Request.Header.FaceCode)
		return nil
	}
	//除了客户端确认的数据包外，其余的包都告诉客户端，已经接收到。
	if c.Request.Header.FaceCode != 200 && c.Request.Header.FaceCode != 520 {
		//告诉客户端已经收到数据包
		c.WriteJsonFaceCode(c.Request.Header.RequestId, 200)
	}
	return c

}
func (s *Server) Close(conn net.Conn) {
	sid := s.GetSessionID(conn)
	//关闭链接
	session, _ := s.MemProvider.SessionRead(sid)
	if session.Conn != nil {
		session.Close()
	}
	<-s.maxClientChan
	//销毁session
	s.MemProvider.SessionDestroy(sid)
}
