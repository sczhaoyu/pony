package srv

import (
	"fmt"
	. "github.com/sczhaoyu/pony/config"
	. "github.com/sczhaoyu/pony/server"
	"github.com/sczhaoyu/pony/util"
)

type ClientServer struct {
	Srv *Server         //链接服务
	CS  *CustomerServer //客户端服务
}

//启动客户端服务
func (c *ClientServer) RunSrv() {
	c.Srv = new(Server)
	c.Srv.Name = "客户端连接服务(1)"
	c.Srv.Handle = c.Handle
	c.Srv.MaxClient = 1
	c.Srv.Run()

}

//链接后台管理
func (c *ClientServer) RunCustomerServer() {
	c.CS = new(CustomerServer)
	c.CS.Name = "客户端后台服务(1)"
	c.CS.ServerAddr = "127.0.0.1:9789"
	c.CS.Handler = HandleCustomer
	c.CS.FirstSend = func() {
		//登录后台服务器管理
		login := &struct {
			Name    string `json:"name"`
			Id      string `json:"id"`
			Addr    string `json:"addr"`
			SrvType int    `json:"srvType"`
		}{
			Name:    c.Srv.Name,
			Id:      c.Srv.Id,
			Addr:    fmt.Sprintf(c.Srv.Ip+":"+"%d", c.Srv.Port),
			SrvType: SERVER_TYPE_CS,
		}
		c.CS.WriteJson(login, 100)
		//获取服务器列表
		c.CS.WriteJson("", 101)
	}
	c.CS.Run()
}
func HandleCustomer(conn *CustomerServer, rsp *Respon) {
	if rsp.Header.FaceCode == 101 {
		ret := []string{}
		conn.Unmarshal(&ret)
		fmt.Println(ret)
	}
}
func (c *ClientServer) Handle(conn *Conn) {
	c.CS.DataChan <- util.ByteLen(conn.Request.Marshal())

}
